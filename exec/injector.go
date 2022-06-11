package exec

import (
	"fmt"

	"github.com/iv-menshenin/hideme/crypt"
	"github.com/iv-menshenin/hideme/message"
)

type InjectConfig interface {
	GetPayload() string
	GetInput() string
	GetOutput() string
	GetPrivateKey() string
	GetAesKey() []byte
	GetSyncKey() []byte
}

func Inject(config InjectConfig) error {
	msg, err := message.NewFromFile(config.GetPayload())
	if err != nil {
		return fmt.Errorf("cannot prepare msg: %w", err)
	}

	carr, err := getCarrier(config.GetInput())
	if err != nil {
		return fmt.Errorf("cannot prepare carrier file: %w", err)
	}
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
	if err = carr.SaveTo(config.GetOutput()); err != nil {
		return fmt.Errorf("cannot save image file: %w", err)
	}
	return nil
}
