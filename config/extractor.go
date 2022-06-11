package config

import (
	"errors"
	"flag"
)

type (
	extractor struct {
		input

		hasSyncKey
		hasAesKey
		hasPublicKey
	}
)

func (e *extractor) initParameters() parser {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	fs.StringVar(&e.input.value, "input", "", "A file that carries hidden information")
	fs.StringVar(&e.publicKey, "public", "", "Public key file path")
	fs.StringVar(&e.syncKeyName, "decode-key", "", "Synchronous key file")
	fs.StringVar(&e.aesKeyName, "aes-key", "", "AES key hex data")
	return fs
}

func (e *extractor) validate() error {
	if e.input.value == "" {
		return errors.New("`input` parameter cannot be empty")
	}
	return nil
}
