package crypto

import (
	"testing"
)

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateKey(-1067048330)
	}
}
