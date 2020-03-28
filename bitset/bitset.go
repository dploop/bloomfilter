package bitset

type Bitset interface {
	Mark(p uint64)
	Test(p uint64) bool
}

type unit = byte

const size = 8

type bitset []unit

func New(m uint64) Bitset {
	n := (m + size - 1) / size
	return make(bitset, n)
}

func (bs bitset) Mark(p uint64) {
	u := unit(1) << (p % size)
	bs[p/size] |= u
}

func (bs bitset) Test(p uint64) bool {
	u := unit(1) << (p % size)
	return bs[p/size]&u == 0
}
