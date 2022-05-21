package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/iv-menshenin/hideme/crypt"
	"log"
	"os"

	"github.com/iv-menshenin/hideme/carrier"
	"github.com/iv-menshenin/hideme/message"
)

type cmd struct {
	command    string
	input      string
	payload    string
	output     string
	privateKey string
	publicKey  string
	syncKey    []byte
	aesKey     []byte
}

const (
	cmdInject   = "inject"
	cmdExtract  = "extract"
	cmdGenerate = "keys"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("available commands: %s, %s or %s", cmdInject, cmdExtract, cmdGenerate)
	}

	var config = parseCmd(os.Args[1], os.Args[2:])
	switch config.command {
	case cmdInject:
		if config.input == "" {
			log.Fatal("`carrier` parameter cannot be empty")
		}
		if config.output == "" {
			log.Fatal("`out` parameter cannot be empty")
		}
		if config.payload == "" {
			log.Fatal("`payload` parameter cannot be empty")
		}
		if err := inject(config); err != nil {
			log.Fatal(err)
		}

	case cmdExtract:
		if config.input == "" {
			log.Fatal("`input` parameter cannot be empty")
		}
		if err := extract(config); err != nil {
			log.Fatal(err)
		}

	case cmdGenerate:
		if err := keysGenerate(config); err != nil {
			log.Fatal(err)
		}

	}
}

func parseCmd(toDo string, args []string) cmd {
	var command = cmd{command: toDo}
	var syncKeyFileName string
	var aesKey string

	switch toDo {
	case cmdInject:
		fs := flag.NewFlagSet(toDo, flag.ExitOnError)
		fs.StringVar(&command.input, "carrier", "", "A PNG file that will carry the valuable information")
		fs.StringVar(&command.payload, "payload", "", "The file you want to hide from prying eyes")
		fs.StringVar(&command.output, "out", "", "The final file, which does not differ from the original file. But it contains encrypted information")
		fs.StringVar(&command.privateKey, "private", "", "Private key file path")
		fs.StringVar(&syncKeyFileName, "encode-key", "", "Synchronous key file")
		fs.StringVar(&aesKey, "aes-key", "", "AES key hex data")
		if err := fs.Parse(args); err != nil {
			log.Fatal(err)
		}

	case cmdExtract:
		fs := flag.NewFlagSet(toDo, flag.ExitOnError)
		fs.StringVar(&command.input, "input", "", "A file that carries hidden information")
		fs.StringVar(&command.publicKey, "public", "", "Public key file path")
		fs.StringVar(&syncKeyFileName, "decode-key", "", "Synchronous key file")
		fs.StringVar(&aesKey, "aes-key", "", "AES key hex data")
		if err := fs.Parse(args); err != nil {
			log.Fatal(err)
		}

	case cmdGenerate:
		fs := flag.NewFlagSet(toDo, flag.ExitOnError)
		fs.StringVar(&command.output, "out", "rsa_key", "Private key file name. `rsa_key` by default.")
		if err := fs.Parse(args); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("available commands: %s, %s\nunknown command: %s", cmdInject, cmdExtract, toDo)
	}

	var err error
	if syncKeyFileName != "" {
		command.syncKey, err = os.ReadFile(syncKeyFileName)
		if err != nil {
			log.Fatalf("cannot open sync file: %s", err)
		}
	}
	if aesKey != "" {
		command.aesKey, err = hex.DecodeString(aesKey)
		if err != nil {
			log.Fatalf("cannot parse aes key as HEX data: %s", err)
		}
	}

	return command
}

const signFileName = "SIGN_FILE"

func inject(config cmd) error {
	msg, err := message.NewFromFile(config.payload)
	if err != nil {
		return fmt.Errorf("cannot prepare msg: %w", err)
	}

	carr, err := getCarrier(config.input)
	if err != nil {
		return fmt.Errorf("cannot prepare carrier file: %w", err)
	}
	secretData := msg.Encode()

	if config.privateKey != "" {
		sign, err := crypt.SignData(secretData, config.privateKey)
		if err != nil {
			return fmt.Errorf("cannot inject secret data to image: %w", err)
		}
		signed, err := message.NewFromBytes(signFileName, sign)
		if err != nil {
			return fmt.Errorf("cannot inject secret data to image: %w", err)
		}
		secretData = append(secretData, signed.Encode()...)
	}

	if len(config.aesKey) > 0 {
		secretData, err = crypt.EncryptDataAES(secretData, config.aesKey)
		if err != nil {
			return fmt.Errorf("cannot encrypt data by aes: %w", err)
		}
	}

	if len(config.syncKey) > 0 {
		err = crypt.EncryptDecryptData(secretData, config.syncKey)
		if err != nil {
			return fmt.Errorf("cannot encode data by key: %w", err)
		}
	}

	if err = carr.Inject(secretData); err != nil {
		return fmt.Errorf("cannot inject secret data to image: %w", err)
	}
	if err = carr.SaveTo(config.output); err != nil {
		return fmt.Errorf("cannot save image file: %w", err)
	}
	return nil
}

func extract(config cmd) error {
	carr, err := getCarrier(config.input)
	if err != nil {
		return fmt.Errorf("cannot prepare carrier file: %w", err)
	}
	data := carr.GetPayload()

	if len(config.syncKey) > 0 {
		err = crypt.EncryptDecryptData(data, config.syncKey)
		if err != nil {
			return fmt.Errorf("cannot encode data by key: %w", err)
		}
	}

	if len(config.aesKey) > 0 {
		data, err = crypt.DecryptDataAES(data, config.aesKey)
		if err != nil {
			return fmt.Errorf("cannot decrypt data by aes: %w", err)
		}
	}

	msgs, err := message.Decode(data)
	if err != nil {
		return fmt.Errorf("cannot decode file from data: %w", err)
	}

	for i, msg := range msgs {
		switch config.publicKey != "" {

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
				if err = crypt.SignVerify(secretData, msg.Content(), config.publicKey); err != nil {
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

type injector interface {
	Inject([]uint8) error
	GetPayload() []uint8
	SaveTo(string) error
}

func getCarrier(fileName string) (injector, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return carrier.New(f)
}

func keysGenerate(config cmd) error {
	private, err := crypt.GenerateKeys()
	if err != nil {
		return fmt.Errorf("cannot generate keys: %w", err)
	}
	if err = crypt.SaveKeysToFile(private, config.output); err != nil {
		return fmt.Errorf("cannot save keys to file `%s`: %w", config.output, err)
	}
	return nil
}
