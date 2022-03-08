package content

const (
	Bit      = 1
	KiloBit  = 1024 * Bit
	MegaBit  = 1 * 1024 * 1024
	Byte     = 8 * Bit
	KiloByte = 1024 * Byte
)

type DataLen uint64

func (dl DataLen) Bits() uint64 {
	return uint64(dl)
}

func (dl DataLen) Bytes() float64 {
	return float64(dl) / Byte
}

func (dl DataLen) KiloBytes() float64 {
	return float64(dl) / KiloByte
}

func (dl DataLen) MegaBites() float64 {
	return float64(dl) / MegaBit
}
