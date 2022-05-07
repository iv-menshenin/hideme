package carrier

import (
	"image/color"
)

func (c *carrier) GetPayload() []uint8 {
	var data []uint8
	for x := 0; x < c.rect.Dx(); x++ {
		for y := 0; y < c.rect.Dy(); y++ {
			col := c.png.At(x, y)
			data = append(data, decodeColor(col))
		}
	}
	return data
}

func decodeColor(c color.Color) uint8 {
	r, g, b, a := c.RGBA()
	r1 := uint8(r)
	g1 := uint8(g)
	b1 := uint8(b)
	a1 := uint8(a)
	switch v := c.(type) {
	case color.NRGBA:
		r1, g1, b1, a1 = v.R, v.G, v.B, v.A
	}
	r1 = r1 & 1
	g1 = (g1 & 1) << 1
	b1 = (b1 & 1) << 2
	a1 = (a1 & 1) << 3
	return r1 | g1 | b1 | a1
}
