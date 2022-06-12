package carrier

import (
	"image/color"

	"github.com/iv-menshenin/hideme/nibbles"
)

func (c *carrier) GetPayload() []uint8 {
	var (
		nn   int
		data = make([]byte, c.extractSz())
	)

xIter:
	for x := 0; x < c.rect.Dx(); x++ {
		if x == xLineReserved {
			continue
		}
		for y := 0; y < c.rect.Dy(); y++ {
			col := c.png.At(x, y)
			data[nn] = decodeColor(col)
			if nn++; nn >= len(data) {
				break xIter
			}
		}
	}
	return nibbles.Convert(data, nibbleSize)
}

func (c *carrier) extractSz() int64 {
	var szChW []byte
	for y := 0; y < c.rect.Dy(); y++ {
		var col = c.png.At(xLineReserved, y)
		szChW = append(szChW, decodeColor(col))
	}
	b := nibbles.Convert(szChW, nibbleSize)
	var sz [8]byte
	copy(sz[:], b)
	return b64int(sz)
}

func b64int(d [8]byte) (result int64) {
	for n := 7; n > 0; n-- {
		result = result | int64(d[n])
		result = result << 8
	}
	return result | int64(d[0])
}

func decodeColor(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	r1 := uint8(r)
	g1 := uint8(g)
	b1 := uint8(b)
	switch v := c.(type) {
	case color.NRGBA:
		r1, g1, b1 = v.R, v.G, v.B
	}
	r1 = r1 & 1
	g1 = (g1 & 1) << 1
	b1 = (b1 & 1) << 2
	return r1 | g1 | b1
}
