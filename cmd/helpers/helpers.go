package helpers

import (
	"crypto/rand"
	"encoding/hex"
)

// CreateRandomHash returns a random hash of length 8
func CreateRandomHash() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}
