package bloomfilter

func bitsetNew(m uint64) []byte {
	return make([]byte, ((m-1)>>3)+1)
}

func bitsetMark(bs []byte, h uint64) {
	b := byte(1) << (h & 7)
	bs[h>>3] |= b
}

func bitsetTest(bs []byte, h uint64) bool {
	b := byte(1) << (h & 7)
	return bs[h>>3]&b == 0
}
