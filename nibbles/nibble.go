package nibbles

type nibble struct {
	mask    int16
	size    int
	current int
	data    []byte
}

const (
	MaxNibbleSize     = 6
	MinNibbleSize     = 1
	DefaultNibbleSize = 4
	bitsInByte        = 8
)

func New(size int, data []byte) *nibble {
	var mask int16
	if size < MinNibbleSize || size > MaxNibbleSize {
		size = DefaultNibbleSize
	}
	for i := 0; i < size; i++ {
		mask |= 1 << i
	}
	return &nibble{
		mask: mask,
		size: size,
		data: data,
	}
}

func (n *nibble) Next() (byte, bool) {
	byteIndex := (n.current * n.size) / bitsInByte
	if byteIndex >= len(n.data) {
		return 0, false
	}
	bitIndex := (n.current * n.size) % bitsInByte
	n.current++
	word := int16(n.data[byteIndex])
	if len(n.data) > byteIndex+1 { // && bitIndex > bitsInByte - n.size
		word |= int16(n.data[byteIndex+1]) << bitsInByte
	}
	result := (word >> bitIndex) & n.mask
	return byte(result), true
}

func Convert(data []byte, size int) (result []byte) {
	var (
		filledBits int
		bitBuffer  int16
	)
	for _, b := range data {
		bitBuffer |= int16(b) << filledBits
		filledBits += size
		if filledBits >= bitsInByte {
			result = append(result, byte(bitBuffer&0xff))
			bitBuffer = bitBuffer >> bitsInByte
			filledBits -= bitsInByte
		}
	}
	if filledBits >= size {
		result = append(result, byte(bitBuffer&0xff))
	}
	return
}
