package carrier

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
)

func (c *carrier) Inject(secret []uint8) error {
	var i = 0
	getMessage := func() (result uint8) {
		if len(secret) > i {
			result = secret[i]
			i++
		}
		return
	}
	for x := 0; x < c.rect.Dx(); x++ {
		for y := 0; y < c.rect.Dy(); y++ {
			col := c.png.At(x, y)
			c.new.Set(x, y, encodeColor(col, getMessage()))
		}
	}
	if len(secret) > i {
		return fmt.Errorf("data lost: %d bytes. choose bigger input file please", len(secret)-i)
	}
	return nil
}

const mask = math.MaxUint8 - 1

func encodeColor(c color.Color, i uint8) color.Color {
	r, g, b, a := c.RGBA()
	r1 := uint8(r>>8)&mask | ((i & 1) >> 0)
	g1 := uint8(g>>8)&mask | ((i & 2) >> 1)
	b1 := uint8(b>>8)&mask | ((i & 4) >> 2)
	a1 := uint8(a>>8)&mask | ((i & 8) >> 3)
	return color.NRGBA{
		R: r1,
		G: g1,
		B: b1,
		A: a1,
	}
}

func (c *carrier) SaveTo(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, c.new)
}
