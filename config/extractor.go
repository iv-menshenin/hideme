package config

import (
	"errors"
	"flag"
	"fmt"
)

type (
	extractor struct {
		input

		hasSyncKey
		hasAesKey
		hasPublicKey
	}
)

func (e *extractor) initCmdParameters() parser {
	fs := flag.NewFlagSet("extract", flag.ExitOnError)
	fs.StringVar(&e.input.value, "input", "", "A file that carries hidden information")
	fs.StringVar(&e.publicKey, "public", "", "Public key file path")
	fs.StringVar(&e.syncKeyName, "decode-key", "", "Synchronous key file")
	fs.StringVar(&e.aesKeyName, "aes-key", "", "AES key hex data")
	return cmdExtractorParser{
		e:      e,
		parser: fs,
	}
}

type cmdExtractorParser struct {
	e      *extractor
	parser parser
}

func (p cmdExtractorParser) Parse(arguments []string) error {
	if err := p.parser.Parse(arguments); err != nil {
		return err
	}
	if err := p.e.hasAesKey.decodeAesKey(); err != nil {
		return fmt.Errorf("can't decode aes key: %s", err)
	}
	if err := p.e.hasSyncKey.loadSyncKey(); err != nil {
		return fmt.Errorf("can't load sync key: %s", err)
	}
	if err := p.e.input.prepare(); err != nil {
		return fmt.Errorf("cannot prepare carrier file: %w", err)
	}
	return nil
}

func (e *extractor) validate() error {
	if e.input.value == "" {
		return errors.New("`input` parameter cannot be empty")
	}
	return nil
}
