package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
)

// CreateRandomHash returns a random hash of length 8
func CreateRandomHash() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

// IsValidURL validates a given URL
func IsValidURL(rawurl string) bool {
	_, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return false
	}
	return true
}
