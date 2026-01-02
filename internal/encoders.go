package internal

import (
	"encoding/binary"
	"fmt"
)

func EncodeChunkHeader(len uint32, nonce []byte) []byte {

	buf := make([]byte, ChunkHeaderSize)

	copy(buf[:12], nonce)
	binary.BigEndian.PutUint32(buf[12:16], len)

	return buf
}

func DecodeChunkHeader(buf []byte) (ChunkHeader, error) {
	var header ChunkHeader

	if len(buf) < ChunkHeaderSize {
		return header, fmt.Errorf("buffer too small, need: %d bytes, got: %d", ChunkHeaderSize, len(buf))
	}

	copy(header.Nonce[:], buf[:12])
	header.Length = binary.BigEndian.Uint32(buf[12:16])

	return header, nil
}

func EncodeHeader(hash []byte, salt []byte) []byte {
	buf := make([]byte, HeaderSize)

	// if len(buf) < int(HeaderSize) {
	// 	return fmt.Errorf("buffer too small, need: %d bytes, got: %d", HeaderSize, len(buf))
	// }

	copy(buf[:4], MagicHeader[:])
	copy(buf[4:20], salt)

	return buf
}

func DecodeHeader(buf []byte) (Header, error) {

	var header Header

	if len(buf) < HeaderSize {
		return header, fmt.Errorf("buffer too small, need: %d bytes, got: %d", HeaderSize, len(buf))
	}

	copy(header.Magic[:], buf[:4])
	copy(header.Salt[:], buf[4:20])

	if header.Magic != MagicHeader {
		return header, fmt.Errorf("invalid magic: expected %q", MagicHeader)
	}

	return header, nil
}
