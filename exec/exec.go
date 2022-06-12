package exec

import (
	"bytes"
	"fmt"
	"os"

	"github.com/iv-menshenin/hideme/carrier"
)

const signFileName = "SIGN_FILE"

type dataCarrier interface {
	FileName() string
	Content() []byte
}

func saveFile(msg dataCarrier) error {
	f, err := os.Create("./" + msg.FileName())
	if err != nil {
		return fmt.Errorf("cannot create new file `%s`: %w", msg.FileName(), err)
	}
	defer f.Close()
	_, err = f.Write(msg.Content())
	if err != nil {
		return fmt.Errorf("cannot write data to file `%s`: %w", msg.FileName(), err)
	}
	return nil

}

type Carrier interface {
	Inject([]uint8) error
	GetPayload() []uint8
	SaveTo(string) error
}

func NewCarrierFromFile(fileName string) (Carrier, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return carrier.New(f)
}

func NewCarrierFromBytes(b []byte) (Carrier, error) {
	return carrier.New(bytes.NewReader(b))
}
