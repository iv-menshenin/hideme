package carrier

import (
	"image/color"
	"math"
	"os"
	"reflect"
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
		{c: color.RGBA{R: 0, G: 0, B: 0, A: math.MaxUint8}, i: 7},
		{c: color.RGBA{R: 128, G: 127, B: 0, A: math.MaxUint8}, i: 6},
		{c: color.RGBA{R: 255, G: 255, B: 255, A: math.MaxUint8}, i: 5},
		{c: color.RGBA{R: 128, G: 127, B: 64, A: math.MaxUint8}, i: 4},
		{c: color.RGBA64{R: 2313, G: 0, B: 257, A: math.MaxUint16}, i: 3},
		{c: color.RGBA64{R: 2313, G: 0, B: 257, A: math.MaxUint16}, i: 2},
		{c: color.Gray{}, i: 1},
		{c: color.White, i: 5},
		{c: color.Black, i: 3},
	}
	for _, c := range testCases {
		enc := encodeColor(c.c, c.i)
		d := decodeColor(enc)
		if c.i != d {
			t.Fatalf("matching error: %d != %d", c.i, d)
		}
	}
}

func TestPNGCarrier(t *testing.T) {
	const tmpFileName = "./test/tmp.png"
	defer os.Remove(tmpFileName)
	var data = []uint8{0, 1, 2, 3, 4, 5, 6, 7, 0, 11, 12, 13, 14, 15, 26, 27, 100, 111, 112, 113, 114, 125, 126, 147, 160, 171, 182, 193, 214, 215, 226, 237, 250, 255}

	makePNGCarrier(t, data, tmpFileName)
	if t.Failed() {
		return
	}
	checkPNGCarrier(t, data, tmpFileName)
}

func makePNGCarrier(t *testing.T, data []uint8, fileName string) {
	f, err := os.Open("./test/test.png")
	if err != nil {
		t.Fatalf("cannot open test.png: %s", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("cannot close opened file: %s", err)
		}
	}()

	img, err := New(f)
	if err != nil {
		t.Errorf("cannot initialize carrier image: %s", err)
		return
	}

	if err = img.Inject(data); err != nil {
		t.Errorf("cannot inject data: %s", err)
		return
	}
	if err = img.SaveTo(fileName); err != nil {
		t.Errorf("cannot save data: %s", err)
		return
	}
}

func checkPNGCarrier(t *testing.T, data []uint8, fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("cannot open test.png: %s", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("cannot close opened file: %s", err)
		}
	}()

	img, err := New(f)
	if err != nil {
		t.Errorf("cannot initialize carrier image: %s", err)
		return
	}

	if err = img.Inject(data); err != nil {
		t.Errorf("cannot inject data: %s", err)
		return
	}
	payload := img.GetPayload()
	if !reflect.DeepEqual(payload, data) {
		t.Errorf("matching error:\nwant: %+v\n got: %+v", data, payload)
	}
}
