package bloomfilter

import (
	"github.com/spaolacci/murmur3"
)

type Hasher func(item []byte) Hashes

type Hashes func(i uint64) uint64

func hasherMurmur3(item []byte) Hashes {
	a, b := murmur3.Sum128(item)
	return func(i uint64) uint64 {
		return a + i*b
	}
}
