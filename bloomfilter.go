package bloomfilter

import (
	"fmt"
	"math"
	"sync"
)

// A BloomFilter is a space-efficient probabilistic data structure, which is
// used to test whether an item is a member of a set. False positive matches
// are possible, but false negatives are not. In other words, a query returns
// either "possibly in set" or "definitely not in set".
type BloomFilter struct {
	// The capacity of bitset.
	m uint64
	// The storage of bitset.
	bitset []byte
	// Number of hash functions.
	k uint64
	// Hasher to generate hashes.
	hasher Hasher
	// A lock to let concurrency.
	lock sync.RWMutex
}

// NewWithEstimate create a Bloom filter for about n items with p false positive
// possibility.
func NewWithEstimate(n uint64, p float64) (*BloomFilter, error) {

	if n == 0 {
		return nil, fmt.Errorf("invalid argument n(%v)", n)
	}

	if p <= 0 || p >= 1 {
		return nil, fmt.Errorf("invalid argument p(%v)", p)
	}

	return New(EstimateParameters(n, p))
}

// EstimateParameters estimates the capacity of bitset and the number of hash
// functions for about n items with p false positive possibility.
func EstimateParameters(n uint64, p float64) (uint64, uint64) {
	m := math.Ceil(math.Log2E * math.Log2(1/p) * float64(n))
	k := math.Ceil(math.Ln2 * m / float64(n))
	return uint64(m), uint64(k)
}

// New create a Bloom filter with m bits storage and k hash functions.
func New(m uint64, k uint64) (*BloomFilter, error) {

	if m == 0 {
		return nil, fmt.Errorf("invalid argument m(%v)", m)
	}

	if k == 0 {
		return nil, fmt.Errorf("invalid argument k(%v)", k)
	}

	bf := &BloomFilter{
		m:      m,
		k:      k,
		hasher: hasherMurmur3,
		bitset: bitsetNew(m),
	}
	return bf, nil
}

// Add operation add item to the Bloom filter.
func (bf *BloomFilter) Add(item []byte) {
	hashes := bf.hasher(item)
	bf.lock.Lock()
	defer bf.lock.Unlock()
	for i := uint64(0); i < bf.k; i++ {
		bitsetMark(bf.bitset, hashes(i)%bf.m)
	}
}

// Contains returns true if the item is in the Bloom filter, false otherwise.
// If true, the result might be a false positive.
// If false, the item is definitely not in the set.
func (bf *BloomFilter) Contains(item []byte) bool {
	hashes := bf.hasher(item)
	bf.lock.RLock()
	defer bf.lock.RUnlock()
	for i := uint64(0); i < bf.k; i++ {
		if bitsetTest(bf.bitset, hashes(i)%bf.m) {
			return false
		}
	}
	return true
}
