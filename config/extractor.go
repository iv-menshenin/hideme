package config

import (
	"errors"
	"flag"
	"fmt"
)

type (
	extractor struct {
		output
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
	fs.StringVar(&e.keyStr, "aes-key", "", "AES key hex data")
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

func (e *extractor) files() []string {
	return e.getFiles()
}

func extractorFromQuery(q Query) (*extractor, error) {
	carr, err := inputFromQuery(q, "carrier")
	if err != nil {
		return nil, fmt.Errorf("can't extract carrier from query: %s", err)
	}
	skey, err := syncKeyFromQuery(q, "encode-key")
	if err != nil {
		return nil, fmt.Errorf("can't extract sync key from query: %s", err)
	}

	e := extractor{
		input:        *carr,
		output:       output{path: "/tmp/hideme/"},
		hasSyncKey:   *skey,
		hasAesKey:    hasAesKey{keyStr: q.StringVal("aes")},
		hasPublicKey: hasPublicKey{publicKey: q.StringVal("public")},
	}
	if err = e.hasAesKey.decodeAesKey(); err != nil {
		return nil, err
	}
	return &e, e.validate()
}
