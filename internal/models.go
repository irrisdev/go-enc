package internal

const (
	HeaderSize      = 20        // bytes
	ChunkHeaderSize = 16        // bytes
	RWSize          = 64 * 1024 // 64 KB
	ChunkSize       = 1 << 20   // 1 MiB
)

var MagicHeader = [4]byte{'g', 'e', 'n', 'c'}

type Header struct {
	Magic [4]byte
	Salt  [16]byte
}

type ChunkHeader struct {
	Nonce  [12]byte
	Length uint32
}
