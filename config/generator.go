package config

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
)

type generator struct {
	output
}

func (g *generator) initCmdParameters() parser {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.StringVar(&g.output.path, "out", "rsa_key", "Private key file name. `rsa_key` by default.")
	return fs
}

func (g *generator) validate() error {
	return nil
}

func (g *generator) files() []string {
	return g.getFiles()
}

func generatorFromQuery(Query) (*generator, error) {
	g := generator{
		output: output{path: makeGeneratedTmpFileName()},
	}
	return &g, g.validate()
}

func makeGeneratedTmpFileName() string {
	var r [16]byte
	if _, err := rand.Read(r[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("/tmp/hideme/%s", hex.EncodeToString(r[:]))
}
