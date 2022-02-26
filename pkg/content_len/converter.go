package content_len

func DataLen(bytesCount float64) *converter {
	return &converter{
		bytesCount: bytesCount,
	}
}

type converter struct {
	bytesCount float64
}

func (c converter) Bites() uint64 {
	return uint64(c.bytesCount * 8)
}

func (c converter) MegaBites() float64 {
	return c.bytesCount * 8 / 1024 / 1024
}
