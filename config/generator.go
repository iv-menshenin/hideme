package config

import (
	"flag"
)

type generator struct {
	output
}

func (g *generator) initParameters() parser {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	fs.StringVar(&g.output.value, "out", "rsa_key", "Private key file name. `rsa_key` by default.")
	return fs
}

func (g *generator) validate() error {
	return nil
}
