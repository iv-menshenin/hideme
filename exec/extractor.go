package exec

import (
	"fmt"
	"log"

	"github.com/iv-menshenin/hideme/crypt"
	"github.com/iv-menshenin/hideme/message"
)

type ExtractConfig interface {
	GetInput() Carrier
	GetPublicKey() string
	GetAesKey() []byte
	GetSyncKey() []byte
}

func Extract(config ExtractConfig) error {
	var err error
	carr := config.GetInput()
	data := carr.GetPayload()

	if syncKey := config.GetSyncKey(); len(syncKey) > 0 {
		err = crypt.EncryptDecryptData(data, syncKey)
		if err != nil {
			return fmt.Errorf("cannot encode data by key: %w", err)
		}
	}

	if aesKey := config.GetAesKey(); len(aesKey) > 0 {
		data, err = crypt.DecryptDataAES(data, aesKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt data by aes: %w", err)
		}
	}

	msgs, err := message.Decode(data)
	if err != nil {
		return fmt.Errorf("cannot decode file from data: %w", err)
	}

	publicKey := config.GetPublicKey()
	for i, msg := range msgs {
		switch publicKey != "" {

		// without sign checking
		case false:
			if msg.FileName() == signFileName {
				// not give away our secret that the file is signed
				continue
			}
			if err = saveFile(&msg); err != nil {
				return fmt.Errorf("cannot save file `%s`: %w", msg.FileName(), err)
			}

		// with sign checking
		case true:
			if msg.FileName() == signFileName {
				if i == 0 {
					return fmt.Errorf("signature cannot be first")
				}
				secretData := msgs[i-1].Encode()
				if err = crypt.SignVerify(secretData, msg.Content(), publicKey); err != nil {
					return fmt.Errorf("cannot verify signature: %w", err)
				}
				log.Println("the signature is verified well")
				if err = saveFile(&msgs[i-1]); err != nil {
					return fmt.Errorf("cannot save file `%s`: %w", msgs[i-1].FileName(), err)
				}
			}

		}
	}
	return nil
}
