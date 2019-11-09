package bloomfilter

import (
	"github.com/spaolacci/murmur3"
)

type Hasher func(data []byte) Hashes

type Hashes func(nth uint64) uint64

func defaultHasher(data []byte) Hashes {
	a, b := murmur3.Sum128(data)
	return func(nth uint64) uint64 {
		return a + nth*b
	}
}
