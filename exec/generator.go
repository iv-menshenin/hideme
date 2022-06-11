package exec

import (
	"fmt"

	"github.com/iv-menshenin/hideme/crypt"
)

type GenerateConfig interface {
	GetOutput() string
}

func Generate(config GenerateConfig) error {
	private, err := crypt.GenerateKeys()
	if err != nil {
		return fmt.Errorf("cannot generate keys: %w", err)
	}
	if err = crypt.SaveKeysToFile(private, config.GetOutput()); err != nil {
		return fmt.Errorf("cannot save keys to file `%s`: %w", config.GetOutput(), err)
	}
	return nil
}
