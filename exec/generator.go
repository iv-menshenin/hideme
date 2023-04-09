package exec

import (
	"fmt"
	"io"

	"github.com/iv-menshenin/hideme/crypt"
)

type GenerateConfig interface {
	SaveFile(string) (io.WriteCloser, error)
}

func Generate(config GenerateConfig) error {
	keys, err := crypt.GenerateKeys()
	if err != nil {
		return fmt.Errorf("cannot generate keys: %w", err)
	}

	wPrivate, err := config.SaveFile("")
	if err != nil {
		return err
	}
	defer wPrivate.Close()

	wPublic, err := config.SaveFile(".pub")
	if err != nil {
		return err
	}
	defer wPublic.Close()

	if err = crypt.SaveKeysToFile(keys, wPublic, wPrivate); err != nil {
		return err
	}
	return nil
}
