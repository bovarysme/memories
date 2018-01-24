package crypto

import (
	"testing"
)

func BenchmarkDeriveKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		deriveKey(-1067048330)
	}
}
