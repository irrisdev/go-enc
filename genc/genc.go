package genc

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"github.com/irrisdev/go-enc/internal"
)

func Encrypt(pass string, filename string) {

	// generate salt
	salt, err := internal.GenerateSalt16()
	if err != nil {
		log.Fatal("failed to generate salt")
	}

	// generate hash using argon2id
	hash, _ := internal.GetArgon2ID(pass, salt)

	// create aes block
	block, err := aes.NewCipher(hash)
	if err != nil {
		log.Fatal("failed to create cipher block")
	}

	// create new gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("failed to create new GCM")

	}

	// open source file
	inFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	// create encrypted destination file - truncates if already exists
	outFile, err := os.Create(fmt.Sprintf("%s.genc", filename))
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// defer recovery function to close and delete incomplete dst file
	defer func() {
		if r := recover(); r != nil {
			log.Printf("closing and removing %s\n", outFile.Name())
			outFile.Close()
			os.Remove(outFile.Name())
			log.Fatalf("recovered panic: %v", r)
		}
	}()

	// create bufio RW buffers limited to RWSize
	reader := bufio.NewReaderSize(inFile, internal.RWSize)
	writer := bufio.NewWriterSize(outFile, internal.RWSize)
	defer writer.Flush()

	// encode and write file header
	header := internal.EncodeHeader(hash, salt)

	// fmt.Printf("writing file header: size = %d bytes\n", len(header))
	writer.Write(header)

	buf := make([]byte, internal.ChunkSize) // 1MiB

	for {
		n, err := reader.Read(buf)
		if n > 0 {

			// generate random nonce
			nonce := make([]byte, gcm.NonceSize())
			_, err = io.ReadFull(rand.Reader, nonce)
			if err != nil {
				panic("failed to generate random nonce")
			}

			// encrypt plaintext using random nonce and buf[:n]
			ciphertext := gcm.Seal(nil, nonce, buf[:n], nil)

			// encode the chunk header
			if len(ciphertext) > math.MaxUint32 {
				panic("slice too large to encode as uint32")
			}

			header := internal.EncodeChunkHeader(uint32(len(ciphertext)), nonce)
			// fmt.Printf("writing chunk header: size = %d bytes\n", len(header))
			// fmt.Printf("cipher size: %d\n", len(ciphertext))

			// write header first then ciphertext
			writer.Write(header)
			writer.Write(ciphertext)

		}

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}

}

func Decrypt(pass string, filename string, outpath ...string) {

	// open file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
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

		log.Println("file already exists, creating a backup")

		if err := internal.CopyFile(path, backup); err != nil {
			log.Fatalf("failed to create backup: %v\n", err)
		}

		log.Printf("backup created at: %s\n", backup)
	}

	outFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	// defer panic function for unexcpected panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("closing %s\n", file.Name())
			file.Close()
			outFile.Close()
			os.Remove(outFile.Name())
			log.Fatalf("recovered panic: %v", r)
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
		log.Fatalf("failed to read file header: %s", err.Error())
	}

	// decode and validate header matches magic
	header, err := internal.DecodeHeader(headerBuf)
	if err != nil {
		log.Fatal(err)
	}

	hash, _ := internal.GetArgon2ID(pass, header.Salt[:])

	// create aes block
	block, err := aes.NewCipher(hash)
	if err != nil {
		log.Fatal("failed to create cipher block")
	}

	// create new gcm
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal("failed to create new GCM")

	}

	for {
		// read 15 byte chunk header
		chunkHeader, err := internal.ReadChunkHeader(reader)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			panic(err)
		}

		// create dynamic buffer to chunk size
		buf := make([]byte, chunkHeader.Length)

		// read full into chunk buffer
		n, err := io.ReadFull(reader, buf)

		if n > 0 {
			// attempt to decrypt chunk using nonce in header
			plaintext, gcmErr := gcm.Open(nil, chunkHeader.Nonce[:], buf, nil)
			if gcmErr != nil {
				panic(gcmErr)
			}
			writer.Write(plaintext)
		}

		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}
	// fmt.Println(string(header.Magic[:]))
	// fmt.Println(base64.RawStdEncoding.EncodeToString(header.Salt[:]))

}
