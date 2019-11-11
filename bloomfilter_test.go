package bloomfilter

import (
	"encoding/binary"
	"runtime"
	"strings"
	"testing"
)

func TestNewWithEstimate(t *testing.T) {

	var (
		bf  *BloomFilter
		err error
	)

	negatives := []struct {
		n uint64
		p float64
		s string
	}{
		{0, 0.3, "invalid argument n"},
		{1000000, -1, "invalid argument p"},
		{1000000, 0, "invalid argument p"},
		{1000000, 1, "invalid argument p"},
		{1000000, 2, "invalid argument p"},
	}
	for _, negative := range negatives {
		bf, err = NewWithEstimate(negative.n, negative.p)
		if bf != nil {
			t.Errorf("bf(%v) should be nil", bf)
		}
		if err == nil || !strings.Contains(err.Error(), negative.s) {
			t.Errorf("err(%v) expected to match %s", err, negative.s)
		}
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
		bf, err = NewWithEstimate(positive.n, positive.p)
		if bf == nil {
			t.Errorf("bf(%v) should not be nil", bf)
		}
		if err != nil {
			t.Errorf("err(%v) should not be nil", err)
		}
	}
}

func TestNew(t *testing.T) {

	var (
		bf  *BloomFilter
		err error
	)

	negatives := []struct {
		m uint64
		k uint64
		s string
	}{
		{0, 5, "invalid argument m"},
		{1000000, 0, "invalid argument k"},
	}
	for _, negative := range negatives {
		bf, err = New(negative.m, negative.k)
		if bf != nil {
			t.Errorf("bf(%v) should be nil", bf)
		}
		if err == nil || !strings.Contains(err.Error(), negative.s) {
			t.Errorf("err(%v) expected to match %s", err, negative.s)
		}
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
		bf, err = New(positive.m, positive.k)
		if bf == nil {
			t.Errorf("bf(%v) should not be nil", bf)
		}
		if err != nil {
			t.Errorf("err(%v) should not be nil", err)
		}
	}
}

func TestBloomFilter_Add(t *testing.T) {

	var (
		bf  *BloomFilter
		err error
	)

	bf, err = NewWithEstimate(1000, 0.03)
	if bf == nil {
		t.Errorf("bf(%v) should not be nil", bf)
	}
	if err != nil {
		t.Errorf("err(%v) should not be nil", err)
	}

	foo, bar := []byte("foo"), []byte("bar")
	bf.Add(foo)
	if !bf.Contains(foo) {
		t.Errorf("bf should contains %s with high possibility", foo)
	}
	if bf.Contains(bar) {
		t.Errorf("bf should not contains %s", bar)
	}
}

func TestBloomFilter_Contains(t *testing.T) {

	var (
		bf  *BloomFilter
		err error
	)

	bf, err = NewWithEstimate(1000, 0.03)
	if bf == nil {
		t.Errorf("bf(%v) should not be nil", bf)
	}
	if err != nil {
		t.Errorf("err(%v) should not be nil", err)
	}

	foo, bar := []byte("foo"), []byte("bar")
	bf.Add(foo)
	if !bf.Contains(foo) {
		t.Errorf("bf should contains %s with high possibility", foo)
	}
	if bf.Contains(bar) {
		t.Errorf("bf should not contains %s", bar)
	}
}

func BenchmarkBloomFilter_Add(b *testing.B) {

	const n = 1000000
	bf, _ := NewWithEstimate(n, 0.03)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		item, v := make([]byte, 4), uint32(i%n)
		binary.BigEndian.PutUint32(item, v)
		bf.Add(item)
	}
	runtime.KeepAlive(bf)
}

func BenchmarkBloomFilter_Contains(b *testing.B) {

	const n = 1000000
	bf, _ := NewWithEstimate(n, 0.03)
	for i := 0; i < n; i++ {
		item, v := make([]byte, 4), uint32(i%n)
		binary.BigEndian.PutUint32(item, v)
		bf.Add(item)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		item, v := make([]byte, 4), uint32(i%n)
		binary.BigEndian.PutUint32(item, v)
		runtime.KeepAlive(bf.Contains(item))
	}
}
