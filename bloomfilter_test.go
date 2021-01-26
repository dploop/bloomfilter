package bloomfilter_test

import (
	"encoding/binary"
	"testing"

	"github.com/dploop/bloomfilter"
	"github.com/dploop/bloomfilter/hasher"
	"github.com/stretchr/testify/assert"
)

func TestNewWithEstimate(t *testing.T) {
	negatives := []struct {
		n   uint64
		p   float64
		msg string
	}{
		{0, 0.3, "invalid argument n"},
		{1000000, -1, "invalid argument p"},
		{1000000, 0, "invalid argument p"},
		{1000000, 1, "invalid argument p"},
		{1000000, 2, "invalid argument p"},
	}
	for _, negative := range negatives {
		bf, err := bloomfilter.NewWithEstimate(negative.n, negative.p)
		assert.Nil(t, bf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), negative.msg)
	}

	positives := []struct {
		n uint64
		p float64
	}{
		{10, 0.03},
		{100, 0.03},
		{10, 0.003},
	}
	for _, positive := range positives {
		bf, err := bloomfilter.NewWithEstimate(positive.n, positive.p,
			bloomfilter.WithHasher(hasher.Murmur3))
		assert.NotNil(t, bf)
		assert.NoError(t, err)
	}
}

func TestNew(t *testing.T) {
	negatives := []struct {
		m   uint64
		k   uint64
		msg string
	}{
		{0, 5, "invalid argument m"},
		{1000000, 0, "invalid argument k"},
	}
	for _, negative := range negatives {
		bf, err := bloomfilter.New(negative.m, negative.k)
		assert.Nil(t, bf)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), negative.msg)
	}

	positives := []struct {
		m uint64
		k uint64
	}{
		{10, 5},
		{100, 5},
		{10, 1},
	}
	for _, positive := range positives {
		bf, err := bloomfilter.New(positive.m, positive.k,
			bloomfilter.WithHasher(hasher.Murmur3))
		assert.NotNil(t, bf)
		assert.NoError(t, err)
	}
}

func TestBloomFilter_Add(t *testing.T) {
	bf, err := bloomfilter.NewWithEstimate(1000, 0.03)
	assert.NotNil(t, bf)
	assert.NoError(t, err)

	foo, bar := []byte("foo"), []byte("bar")
	assert.False(t, bf.Contains(foo))
	assert.False(t, bf.Contains(bar))
	bf.Add(foo)
	assert.True(t, bf.Contains(foo))
	assert.False(t, bf.Contains(bar))
	bf.Add(bar)
	assert.True(t, bf.Contains(foo))
	assert.True(t, bf.Contains(bar))
}

func TestBloomFilter_Contains(t *testing.T) {
	bf, err := bloomfilter.NewWithEstimate(1000, 0.03)
	assert.NotNil(t, bf)
	assert.NoError(t, err)

	foo, bar := []byte("foo"), []byte("bar")
	assert.False(t, bf.Contains(foo))
	assert.False(t, bf.Contains(bar))
	bf.Add(foo)
	assert.True(t, bf.Contains(foo))
	assert.False(t, bf.Contains(bar))
	bf.Add(bar)
	assert.True(t, bf.Contains(foo))
	assert.True(t, bf.Contains(bar))
}

func BenchmarkBloomFilter_Add(b *testing.B) {
	const n = 1000000
	bf, _ := bloomfilter.NewWithEstimate(n, 0.03)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		item, v := make([]byte, 8), uint64(i%n)
		binary.BigEndian.PutUint64(item, v)
		bf.Add(item)
	}
}

func BenchmarkBloomFilter_Contains(b *testing.B) {
	const n = 1000000
	bf, _ := bloomfilter.NewWithEstimate(n, 0.03)
	for i := 0; i < n; i++ {
		item, v := make([]byte, 8), uint64(i%n)
		binary.BigEndian.PutUint64(item, v)
		bf.Add(item)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		item, v := make([]byte, 8), uint64(i%n)
		binary.BigEndian.PutUint64(item, v)
		_ = bf.Contains(item)
	}
}
