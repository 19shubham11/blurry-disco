package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRandomHash(t *testing.T) {
	t.Run("Should return string of length 8", func(t *testing.T) {
		got := CreateRandomHash()
		assert.Equal(t, len(got), 8)
	})

	t.Run("Should create 10k unique hashes", func(t *testing.T) {
		set := make(map[string]struct{})
		tenK := 10000
		for i := 0; i < tenK; i++ {
			hash := CreateRandomHash()
			set[hash] = struct{}{}
		}
		assert.Equal(t, len(set), tenK)
	})
}

func BenchmarkCreateRadomHash(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CreateRandomHash()
	}
}
