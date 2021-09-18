package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUniqueHash(t *testing.T) {
	t.Run("Should return string of length 8", func(t *testing.T) {
		got := CreateUniqueHash()
		assert.Equal(t, len(got), 8)
	})

	t.Run("Should create 10k unique hashes", func(t *testing.T) {
		set := make(map[string]struct{})
		tenK := 10000
		for i := 0; i < tenK; i++ {
			hash := CreateUniqueHash()
			set[hash] = struct{}{}
		}
		assert.Equal(t, len(set), tenK)
	})
}

func TestIsValidURL(t *testing.T) {
	t.Run("Should return true for a valid URL", func(t *testing.T) {
		url := "http://www.google.com"
		assert.True(t, IsValidURL(url))
	})

	t.Run("Should return false for a valid URL", func(t *testing.T) {
		url := "notAURL"
		assert.False(t, IsValidURL(url))
	})
}

func BenchmarkCreateRadomHash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CreateUniqueHash()
	}
}
