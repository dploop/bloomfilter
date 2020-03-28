package bloomfilter

import (
	"fmt"
	"math"
	"sync"

	"github.com/dploop/bloomfilter/bitset"
	"github.com/dploop/bloomfilter/hasher"
)

// Bloom filter is a space-efficient probabilistic data structure, which is
// used to test whether an item is a member of a set. False positive matches
// are possible, but false negatives are not. In other words, a query returns
// either "possibly in set" or "definitely not in set".

// The interface of a Bloom filter.
type BloomFilter interface {
	Add(item []byte)
	AddWithoutLock(item []byte)
	Contains(item []byte) bool
	ContainsWithoutLock(item []byte) bool
}

// The implementation of a Bloom filter.
type bloomFilter struct {
	// A mutex to let concurrency.
	sync.RWMutex
	// The capacity of bitset.
	m uint64
	// The storage of bitset.
	bitset.Bitset
	// Number of hash functions.
	k uint64
	// Hasher to generate hashes.
	hasher.Hasher
}

// The option when creating a Bloom filter.
type Option func(*bloomFilter)

// WithHasher creates an option of hasher.
func WithHasher(hs hasher.Hasher) Option {
	return func(bf *bloomFilter) {
		bf.Hasher = hs
	}
}

// NewWithEstimate creates a Bloom filter for about `n` items with `p` false
// positive possibility.
func NewWithEstimate(n uint64, p float64, opts ...Option) (BloomFilter, error) {
	if n == 0 {
		return nil, fmt.Errorf("invalid argument n(%v)", n)
	}

	if p <= 0 || p >= 1 {
		return nil, fmt.Errorf("invalid argument p(%v)", p)
	}

	m, k := EstimateParameters(n, p)

	return New(m, k, opts...)
}

// EstimateParameters estimates the capacity of bitset `m` and the number of
// hash functions `k` for about `n` items with `p` false positive possibility.
func EstimateParameters(n uint64, p float64) (uint64, uint64) {
	m := math.Ceil(math.Log2E * math.Log2(1/p) * float64(n))
	k := math.Ceil(math.Ln2 * m / float64(n))

	return uint64(m), uint64(k)
}

// New creates a Bloom filter with `m` bits storage and `k` hash functions.
func New(m uint64, k uint64, opts ...Option) (BloomFilter, error) {
	if m == 0 {
		return nil, fmt.Errorf("invalid argument m(%v)", m)
	}

	if k == 0 {
		return nil, fmt.Errorf("invalid argument k(%v)", k)
	}

	bf := &bloomFilter{
		m:      m,
		Bitset: bitset.New(m),
		k:      k,
		Hasher: hasher.Murmur3,
	}

	for _, opt := range opts {
		opt(bf)
	}

	return bf, nil
}

// Add adds item to the Bloom filter.
func (bf *bloomFilter) Add(item []byte) {
	bf.Lock()
	defer bf.Unlock()

	bf.AddWithoutLock(item)
}

// AddWithoutLock is same with Add, but without lock.
func (bf *bloomFilter) AddWithoutLock(item []byte) {
	a, b := bf.Hasher(item)
	a, b = a%bf.m, b%bf.m

	for i := uint64(0); i < bf.k; i++ {
		bf.Bitset.Mark((a + i*b) % bf.m)
	}
}

// Contains returns true if the item is in the Bloom filter, false otherwise.
// If true, the result might be a false positive.
// If false, the item is definitely not in the set.
func (bf *bloomFilter) Contains(item []byte) bool {
	bf.RLock()
	defer bf.RUnlock()

	return bf.ContainsWithoutLock(item)
}

// ContainsWithoutLock is same with Contains, but without lock.
func (bf *bloomFilter) ContainsWithoutLock(item []byte) bool {
	a, b := bf.Hasher(item)
	a, b = a%bf.m, b%bf.m

	for i := uint64(0); i < bf.k; i++ {
		if bf.Bitset.Test((a + i*b) % bf.m) {
			return false
		}
	}

	return true
}
