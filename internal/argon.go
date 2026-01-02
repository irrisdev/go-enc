package internal

import (
	"encoding/base64"
	"fmt"
	"runtime"

	"golang.org/x/crypto/argon2"
)

const (
	mem        uint32 = 64 * 1024 // 64 MB
	keyLen     uint32 = 32        // 32 bits
	iterations uint32 = 1
)

func GetArgon2ID(pass string, salt []byte) (key []byte, hash string) {
	// get cpu threads available
	threads := uint8(runtime.NumCPU())
	if threads > 4 {
		threads = 4
	}

	key = argon2.IDKey([]byte(pass), salt, iterations, mem, threads, keyLen)

	// base 64 for sprintf
	bSalt := base64.RawStdEncoding.EncodeToString(salt)
	bKey := base64.RawStdEncoding.EncodeToString(key)

	hash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, mem, iterations, threads, bSalt, bKey)

	return key, hash
}
