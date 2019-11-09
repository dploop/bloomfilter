package bloomfilter

import (
	"fmt"
	"math"
	"unsafe"
)

type bucket = uint

const width = uint64(unsafe.Sizeof(bucket(0)))

type BloomFilter struct {
	m       uint64
	k       uint64
	hasher  Hasher
	buckets []bucket
}

func NewWithEstimate(n uint64, p float64) (*BloomFilter, error) {

	if n == 0 {
		return nil, fmt.Errorf("invalid argument n(%v)", n)
	}

	if p <= 0 || p >= 1 {
		return nil, fmt.Errorf("invalid argument p(%v)", p)
	}

	return New(EstimateParameters(n, p))
}

func EstimateParameters(n uint64, p float64) (uint64, uint64) {
	m := math.Ceil(math.Log2E * math.Log2(1/p) * float64(n))
	k := math.Ceil(math.Ln2 * m / float64(n))
	return uint64(m), uint64(k)
}

func New(m uint64, k uint64) (*BloomFilter, error) {

	if m == 0 {
		return nil, fmt.Errorf("invalid argument m(%v)", m)
	}

	if k == 0 {
		return nil, fmt.Errorf("invalid argument k(%v)", k)
	}

	bf := &BloomFilter{m: m, k: k, hasher: defaultHasher}
	bf.buckets = make([]bucket, (m-1)/width+1)

	return bf, nil
}

func (bf *BloomFilter) Add(data []byte) {
	hashes := bf.hasher(data)
	for i := uint64(0); i < bf.k; i++ {
		h := hashes(i) % bf.m
		b := bucket(1) << (i % width)
		bf.buckets[h/width] |= b
	}
}

func (bf *BloomFilter) Test(data []byte) bool {
	hashes := bf.hasher(data)
	for i := uint64(0); i < bf.k; i++ {
		h := hashes(i) % bf.m
		b := bucket(1) << (i % width)
		if bf.buckets[h/width]&b == 0 {
			return false
		}
	}
	return true
}
