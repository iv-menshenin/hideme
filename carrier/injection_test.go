package carrier

import (
	"image/color"
	"math"
	"testing"
)

func Test_encodeColor(t *testing.T) {
	t.Parallel()
	type test struct {
		c color.Color
		i uint8
	}
	testCases := []test{
		{c: color.RGBA{R: 0, G: 12, B: 3, A: math.MaxUint8}, i: 0},
		{c: color.RGBA{R: 0, G: 0, B: 0, A: math.MaxUint8}, i: 8},
		{c: color.RGBA{R: 128, G: 127, B: 0, A: math.MaxUint8}, i: 12},
		{c: color.RGBA{R: 255, G: 255, B: 255, A: math.MaxUint8}, i: 3},
		{c: color.RGBA{R: 128, G: 127, B: 64, A: math.MaxUint8}, i: 15},
		{c: color.RGBA64{R: 2313, G: 0, B: 257, A: math.MaxUint16}, i: 0},
		{c: color.RGBA64{R: 2313, G: 0, B: 257, A: math.MaxUint16}, i: 8},
		{c: color.Gray{}, i: 15},
		{c: color.White, i: 15},
		{c: color.Black, i: 9},
	}
	for _, c := range testCases {
		enc := encodeColor(c.c, c.i)
		d := decodeColor(enc)
		if c.i != d {
			t.Fatalf("matching error: %d != %d", c.i, d)
		}
	}
}
