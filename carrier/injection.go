package carrier

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/iv-menshenin/hideme/nibbles"
)

const (
	nibbleSize    = 3
	xLineReserved = 0
)

func (c *carrier) Inject(secret []uint8) error {
	var nib = nibbles.New(nibbleSize, secret)

	wrote, err := c.injectData(nib)
	if err != nil {
		return err
	}

	return c.injectDataSize(wrote)
}

type nibbler interface {
	Next() (byte, bool)
}

func (c *carrier) injectData(nib nibbler) (int64, error) {
	var (
		done     bool
		injected int64
	)

	for x := 0; x < c.rect.Dx(); x++ {
		if x == xLineReserved {
			continue
		}
		for y := 0; y < c.rect.Dy(); y++ {
			var col = c.png.At(x, y)
			if !done {
				if b, ok := nib.Next(); ok {
					injected++
					col = encodeColor(col, b)
				} else {
					done = true
				}
			}
			c.new.Set(x, y, col)
		}
	}

	if !done {
		var left int64
		for {
			if _, ok := nib.Next(); !ok {
				break
			}
			left++
		}
		return injected, fmt.Errorf("data lost: %d bytes injected, %d bytes left. choose bigger input file please", injected, left)
	}
	return injected, nil
}

func (c *carrier) injectDataSize(sz int64) error {
	var (
		done bool
		szCh = int64b(sz)
		nib  = nibbles.New(nibbleSize, szCh[:])
	)
	for y := 0; y < c.rect.Dy(); y++ {
		var col = c.png.At(xLineReserved, y)
		if !done {
			if b, ok := nib.Next(); ok {
				col = encodeColor(col, b)
			} else {
				done = true
			}
		}
		c.new.Set(xLineReserved, y, col)
	}
	return nil
}

const mask = math.MaxUint8 - 1

func encodeColor(c color.Color, i uint8) color.Color {
	r, g, b, a := c.RGBA()
	r1 := uint8(r>>8)&mask | ((i & 1) >> 0)
	g1 := uint8(g>>8)&mask | ((i & 2) >> 1)
	b1 := uint8(b>>8)&mask | ((i & 4) >> 2)
	a1 := uint8(a >> 8)
	return color.NRGBA{
		R: r1,
		G: g1,
		B: b1,
		A: a1,
	}
}

func int64b(i int64) (result [8]byte) {
	for n := 0; n < len(result); n++ {
		m := i & math.MaxUint8
		result[n] = byte(m)
		i = i >> 8
	}
	return
}

func (c *carrier) SaveTo(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, c.new)
}
