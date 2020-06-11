package gas

//用一个只读的[]byte切片包装真实缓存值

type ByteView struct {
	b []byte
}

//返回缓存大小
func (v ByteView) Len() int {
	return len(v.b)
}

//返回切片的拷贝，防止修改切片
func (v ByteView) ByteSlice() []byte {
	return byteClone(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func byteClone(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
