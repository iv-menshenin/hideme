package carrier

import (
	"bytes"
	"io"
	"os"
)

type Carrier interface {
	Inject([]uint8) error
	GetPayload() []uint8
	SaveTo(f io.Writer) error
}

func NewCarrierFromFile(fileName string) (Carrier, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return New(f)
}

func NewCarrierFromBytes(b []byte) (Carrier, error) {
	return New(bytes.NewReader(b))
}
