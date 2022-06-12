package config

import (
	"io"
	"net/http"

	"github.com/iv-menshenin/hideme/exec"
)

type (
	Config struct {
		argsKeeper
		args []string
		doer func() error
	}
	argsKeeper interface {
		initCmdParameters() parser
		validate() error
	}
	parser interface {
		Parse(arguments []string) error
	}
	Query interface {
		StringVal(string) string
		ByteVal(string) (io.ReadCloser, string, error)
	}
)

func NewInjector(args []string) *Config {
	var i = &injector{}
	return &Config{
		argsKeeper: i,
		args:       args,
		doer: func() error {
			return exec.Inject(i)
		},
	}
}

func NewInjectorFromQuery(q Query) (*Config, error) {
	i, err := injectorFromQuery(q)
	if err != nil {
		return nil, err
	}
	return &Config{
		argsKeeper: i,
		doer: func() error {
			return exec.Inject(i)
		},
	}, nil
}

func NewExtractor(args []string) *Config {
	var e = &extractor{}
	return &Config{
		argsKeeper: e,
		args:       args,
		doer: func() error {
			return exec.Extract(e)
		},
	}
}

func NewGenerator(args []string) *Config {
	var g = &generator{}
	return &Config{
		argsKeeper: g,
		args:       args,
		doer: func() error {
			return exec.Generate(g)
		},
	}
}

func NewServer(args []string, hf http.HandlerFunc) *Config {
	var s = &server{}
	return &Config{
		argsKeeper: s,
		args:       args,
		doer: func() error {
			return exec.Serve(s, hf)
		},
	}
}

func (c *Config) Parse() error {
	return c.argsKeeper.initCmdParameters().Parse(c.args)
}

func (c *Config) Validate() error {
	return c.argsKeeper.validate()
}

func (c *Config) Execute() error {
	return c.doer()
}
