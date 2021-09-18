package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"
)

// CreateUniqueHash returns a random hash of length 8
func CreateUniqueHash() string {
	byteLength := 4
	bytes := make([]byte, byteLength)
	_, _ = rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

// IsValidURL validates a given URL
func IsValidURL(rawurl string) bool {
	_, err := url.ParseRequestURI(rawurl)

	return err == nil
}
