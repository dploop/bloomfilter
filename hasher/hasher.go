package hasher

import (
	"github.com/spaolacci/murmur3"
)

type Hasher func(item []byte) (uint64, uint64)

var Murmur3 = murmur3.Sum128
