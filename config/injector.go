package config

import (
	"errors"
	"flag"
)

type (
	injector struct {
		payload
		input
		output

		hasSyncKey
		hasAesKey
		hasPrivateKey
	}
)

func (i *injector) initParameters() parser {
	fs := flag.NewFlagSet("inject", flag.ExitOnError)
	fs.StringVar(&i.input.value, "carrier", "", "A PNG file that will carry the valuable information")
	fs.StringVar(&i.payload.value, "payload", "", "The file you want to hide from prying eyes")
	fs.StringVar(&i.output.value, "out", "", "The final file, which does not differ from the original file. But it contains encrypted information")
	fs.StringVar(&i.privateKey, "private", "", "Private key file path")
	fs.StringVar(&i.syncKeyName, "encode-key", "", "Synchronous key file")
	fs.StringVar(&i.aesKeyName, "aes-key", "", "AES key hex data")
	return fs
}

func (i *injector) validate() error {
	if i.input.value == "" {
		return errors.New("`carrier` parameter cannot be empty")
	}
	if i.output.value == "" {
		return errors.New("`out` parameter cannot be empty")
	}
	if i.payload.value == "" {
		return errors.New("`payload` parameter cannot be empty")
	}
	return nil
}
