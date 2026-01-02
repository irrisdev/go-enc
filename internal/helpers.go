package internal

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"io"
	"os"
)

func ReadChunkHeader(reader *bufio.Reader) (ChunkHeader, error) {
	var header ChunkHeader
	// create buffer of size chunkheader
	buf := make([]byte, ChunkHeaderSize)

	// read full bytes into buffer
	n, err := io.ReadFull(reader, buf)

	// nothing left to read, EOF
	if n == 0 && err == io.EOF {
		return header, io.EOF
	}
	// partial read == error
	if n < ChunkHeaderSize && err != nil {
		return header, io.ErrUnexpectedEOF
	}

	// decode header correctly
	header, headerErr := DecodeChunkHeader(buf)
	if headerErr != nil {
		return header, headerErr
	}

	return header, nil
}

func GenerateSalt16() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func GenerateSalt32() ([]byte, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func SaltToString(salt []byte) string {
	return base64.RawStdEncoding.EncodeToString(salt)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}
