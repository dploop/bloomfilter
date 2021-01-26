package hasher

import (
	"github.com/spaolacci/murmur3"
)

type Hasher func(item []byte) (uint64, uint64)

func Murmur3(item []byte) (uint64, uint64) {
	return murmur3.Sum128(item)
}
