package mircool_cache

type CacheData struct {
	b []byte
}

func (c CacheData) Len() int {
	return len(c.b)
}

func (c CacheData) ByteSlice() []byte {
	return cloneBytes(c.b)
}

func (c CacheData) String() string {
	return string(c.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
