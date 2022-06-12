package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
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

func (i *injector) initCmdParameters() parser {
	fs := flag.NewFlagSet("inject", flag.ExitOnError)
	fs.StringVar(&i.input.value, "carrier", "", "A PNG file that will carry the valuable information")
	fs.StringVar(&i.payload.fileName, "payload", "", "The file you want to hide from prying eyes")
	fs.StringVar(&i.output.value, "out", "", "The final file, which does not differ from the original file. But it contains encrypted information")
	fs.StringVar(&i.privateKey, "private", "", "Private key file path")
	fs.StringVar(&i.syncKeyName, "encode-key", "", "Synchronous key file")
	fs.StringVar(&i.keyStr, "aes-key", "", "AES key hex data")
	return cmdInjectorParser{
		i:      i,
		parser: fs,
	}
}

type cmdInjectorParser struct {
	i      *injector
	parser parser
}

func (p cmdInjectorParser) Parse(arguments []string) error {
	if err := p.parser.Parse(arguments); err != nil {
		return err
	}
	if err := p.i.hasAesKey.decodeAesKey(); err != nil {
		return fmt.Errorf("can't decode aes key: %s", err)
	}
	if err := p.i.hasSyncKey.loadSyncKey(); err != nil {
		return fmt.Errorf("can't load sync key: %s", err)
	}
	if err := p.i.input.prepare(); err != nil {
		return fmt.Errorf("cannot prepare carrier file: %w", err)
	}
	if err := p.i.payload.prepare(); err != nil {
		return fmt.Errorf("cannot prepare msg: %w", err)
	}
	return nil
}

func (i *injector) validate() error {
	if i.input.value == "" {
		return errors.New("`carrier` parameter cannot be empty")
	}
	if i.output.value == "" {
		return errors.New("`out` parameter cannot be empty")
	}
	if i.payload.message == nil {
		return errors.New("`payload` parameter cannot be empty")
	}
	return nil
}

func injectorFromQuery(q Query) (*injector, error) {
	carr, err := inputFromQuery(q, "carrier")
	if err != nil {
		return nil, fmt.Errorf("can't extract carrier from query: %s", err)
	}
	pload, err := payloadFromQuery(q, "payload")
	if err != nil {
		return nil, fmt.Errorf("can't extract payload from query: %s", err)
	}
	skey, err := syncKeyFromQuery(q, "encode-key")
	if err != nil {
		return nil, fmt.Errorf("can't extract sync key from query: %s", err)
	}

	i := injector{
		payload:       *pload,
		input:         *carr,
		output:        output{value: makeTmpFileName()},
		hasSyncKey:    *skey,
		hasAesKey:     hasAesKey{keyStr: q.StringVal("aes")},
		hasPrivateKey: hasPrivateKey{privateKey: q.StringVal("private")},
	}
	if err = i.hasAesKey.decodeAesKey(); err != nil {
		return nil, err
	}
	return &i, i.validate()
}

func makeTmpFileName() string {
	var r [16]byte
	if _, err := rand.Read(r[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("/tmp/hideme/%s.png", hex.EncodeToString(r[:]))
}
