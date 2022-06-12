package exec

import (
	"fmt"
	"io"

	"github.com/iv-menshenin/hideme/carrier"
	"github.com/iv-menshenin/hideme/crypt"
	"github.com/iv-menshenin/hideme/message"
)

type InjectConfig interface {
	GetPayload() *message.Message
	GetInput() carrier.Carrier
	GetPrivateKey() string
	GetAesKey() []byte
	GetSyncKey() []byte
	SaveFile(string) (io.WriteCloser, error)
}

func Inject(config InjectConfig) (err error) {
	msg := config.GetPayload()
	carr := config.GetInput()
	secretData := msg.Encode()

	if privateKey := config.GetPrivateKey(); privateKey != "" {
		sign, err := crypt.SignData(secretData, privateKey)
		if err != nil {
			return fmt.Errorf("cannot inject secret data to image: %w", err)
		}
		signed, err := message.NewFromBytes(signFileName, sign)
		if err != nil {
			return fmt.Errorf("cannot inject secret data to image: %w", err)
		}
		secretData = append(secretData, signed.Encode()...)
	}

	if aesKey := config.GetAesKey(); len(aesKey) > 0 {
		secretData, err = crypt.EncryptDataAES(secretData, aesKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt data by aes: %w", err)
		}
	}

	if syncKey := config.GetSyncKey(); len(syncKey) > 0 {
		err = crypt.EncryptDecryptData(secretData, syncKey)
		if err != nil {
			return fmt.Errorf("cannot encode data by key: %w", err)
		}
	}

	if err = carr.Inject(secretData); err != nil {
		return fmt.Errorf("cannot inject secret data to image: %w", err)
	}

	w, err := config.SaveFile("")
	if err != nil {
		return fmt.Errorf("cannot save image file: %w", err)
	}
	defer w.Close()
	if err = carr.SaveTo(w); err != nil {
		return fmt.Errorf("cannot save image file: %w", err)
	}
	return nil
}
