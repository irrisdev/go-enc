package genc

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/irrisdev/go-enc/internal"
)

var (
	ErrNewSalt       = errors.New("failed to generate salt")
	ErrNewCipher     = errors.New("failed to create cipher block")
	ErrNewGcm        = errors.New("failed to create new GCM")
	ErrNewNonce      = errors.New("failed to generate random nonce")
	ErrChunkTooLarge = errors.New("slice too large to encode as uint32")
	ErrRemoveOrigin  = errors.New("failed to remove original file")
	ErrBakFile       = errors.New("failed to backup file")
	ErrSyncEncFile   = errors.New("failed to sync encrypted file")
	ErrOpenFile      = errors.New("failed to open file")
	ErrCreateFile    = errors.New("failed to create file")
	ErrWriteHeader   = errors.New("failed to write file header")
	ErrReadHeader    = errors.New("failed to read file header")
)

func Encrypt(pass string, filename string, deleteOrignal ...bool) error {

	// generate salt
	salt, err := internal.GenerateSalt16()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNewSalt, err)
	}

	// generate hash using argon2id
	hash, _ := internal.GetArgon2ID(pass, salt)

	// create aes block
	block, err := aes.NewCipher(hash)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNewCipher, err)
	}

	// create new gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNewGcm, err)
	}

	// open source file
	inFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOpenFile, err)
	}
	defer inFile.Close()

	// create encrypted destination file - truncates if already exists
	outFile, err := os.Create(fmt.Sprintf("%s.genc", filename))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateFile, err)
	}
	defer outFile.Close()

	// check if is completed when function returns otherwise delete outFile
	completed := false
	defer func() {
		if !completed {
			os.Remove(outFile.Name())
			log.Println("encryption failed")
		}
	}()

	// create bufio RW buffers limited to RWSize
	reader := bufio.NewReaderSize(inFile, internal.RWSize)
	writer := bufio.NewWriterSize(outFile, internal.RWSize)
	defer writer.Flush()

	// encode and write file header
	header := internal.EncodeHeader(hash, salt)

	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("%w: %w", ErrWriteHeader, err)
	}

	buf := make([]byte, internal.ChunkSize) // 1MiB

	for {
		n, err := reader.Read(buf)
		if n > 0 {

			// generate random nonce
			nonce := make([]byte, gcm.NonceSize())
			_, err = io.ReadFull(rand.Reader, nonce)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrNewNonce, err)
			}

			// encrypt plaintext using random nonce and buf[:n]
			ciphertext := gcm.Seal(nil, nonce, buf[:n], nil)

			// encode the chunk header
			if len(ciphertext) > math.MaxUint32 {
				return fmt.Errorf("%w: %w", ErrChunkTooLarge, err)
			}

			header := internal.EncodeChunkHeader(uint32(len(ciphertext)), nonce)

			// write header first then ciphertext
			writer.Write(header)
			writer.Write(ciphertext)

		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	// fsync guarantee file write
	if err := outFile.Sync(); err != nil {
		return fmt.Errorf("%w: %w", ErrSyncEncFile, err)
	}

	completed = true

	if len(deleteOrignal) > 0 && deleteOrignal[0] {
		if err := os.Remove(filename); err != nil {
			return ErrRemoveOrigin
		}
	}

	return nil
}

func Decrypt(pass string, filename string, outpath ...string) error {
	completed := false

	// open file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOpenFile, err)
	}
	defer file.Close()

	// check if outpath has been specified
	var path string
	if len(outpath) > 0 {
		path = outpath[0]
	} else {
		path = strings.TrimSuffix(filename, ".genc")
	}

	// check if file already exists, create backup if it does
	if internal.FileExists(path) {
		backup := fmt.Sprintf("%s.bak", path)

		log.Printf("output file %s already exists, creating backup at: %s\n", filepath.Base(path), backup)

		if err := internal.CopyFile(path, backup); err != nil {
			return fmt.Errorf("%w: %w", ErrBakFile, err)
		}
	}

	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateFile, err)
	}

	defer func() {
		if !completed {
			os.Remove(outFile.Name())
			log.Println("decryption failed")
		}
	}()

	// create buffered io reader
	reader := bufio.NewReaderSize(file, internal.RWSize)
	writer := bufio.NewWriterSize(outFile, internal.RWSize)
	defer writer.Flush()

	// create buffer of exact header size
	headerBuf := make([]byte, internal.HeaderSize)

	// read full 52 bytes of header, err if cannot
	_, err = io.ReadFull(reader, headerBuf)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadHeader, err)
	}

	// decode and validate header matches magic
	header, err := internal.DecodeHeader(headerBuf)
	if err != nil {
		return err
	}

	hash, _ := internal.GetArgon2ID(pass, header.Salt[:])

	// create aes block
	block, err := aes.NewCipher(hash)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNewCipher, err)
	}

	// create new gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrNewGcm, err)
	}

	for {
		// read 15 byte chunk header
		chunkHeader, chunkErr := internal.ReadChunkHeader(reader)
		if chunkErr == io.EOF || chunkErr == io.ErrUnexpectedEOF { // unsure if ErrUnexpectedEOF should return error
			break
		}

		if chunkErr != nil {
			return err
		}

		// create dynamic buffer to chunk size
		buf := make([]byte, chunkHeader.Length)

		// read full into chunk buffer
		n, err := io.ReadFull(reader, buf)

		if n > 0 {
			// attempt to decrypt chunk using nonce in header
			plaintext, gcmErr := gcm.Open(nil, chunkHeader.Nonce[:], buf, nil)
			if gcmErr != nil {
				return fmt.Errorf("%w: %w", gcmErr, gcmErr)
			}
			writer.Write(plaintext)
		}

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			return err
		}
	}

	completed = true

	return nil
}
