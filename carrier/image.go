package carrier

import (
	"image"
	"image/png"
	"io"
)

type carrier struct {
	rect image.Rectangle
	png  image.Image
	new  *image.NRGBA
}

func New(r io.Reader) (*carrier, error) {
	imagePNG, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	var c = carrier{
		rect: imagePNG.Bounds(),
		png:  imagePNG,
	}
	c.new = image.NewNRGBA(c.rect)
	return &c, nil
}
