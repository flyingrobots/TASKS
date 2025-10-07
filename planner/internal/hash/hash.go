package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashCanonicalBytes computes the SHA-256 hash of the input bytes.
// The spec requires hashing the canonical bytes plus a final newline.
func HashCanonicalBytes(input []byte) string {
	hasher := sha256.New()
	hasher.Write(input)
	return hex.EncodeToString(hasher.Sum(nil))
}
