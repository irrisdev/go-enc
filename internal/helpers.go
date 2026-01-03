/*
Copyright © 2026 irrisdev lithium8260@proton.me

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

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
