package config

import "github.com/iv-menshenin/hideme/exec"

type (
	Config struct {
		argsKeeper
		doer func() error
	}
	argsKeeper interface {
		initCmdParameters() parser
		validate() error
	}
	parser interface {
		Parse(arguments []string) error
	}
)

func NewInjector() *Config {
	var i = &injector{}
	return &Config{
		argsKeeper: i,
		doer: func() error {
			return exec.Inject(i)
		},
	}
}

func NewExtractor() *Config {
	var e = &extractor{}
	return &Config{
		argsKeeper: e,
		doer: func() error {
			return exec.Extract(e)
		},
	}
}

func NewGenerator() *Config {
	var g = &generator{}
	return &Config{
		argsKeeper: g,
		doer: func() error {
			return exec.Generate(g)
		},
	}
}

func NewServer() *Config {
	var s = &server{}
	return &Config{
		argsKeeper: s,
		doer: func() error {
			return exec.Serve(s)
		},
	}
}

func (c *Config) Parse(args []string) error {
	return c.argsKeeper.initCmdParameters().Parse(args)
}

func (c *Config) Validate() error {
	return c.argsKeeper.validate()
}

func (c *Config) Execute() error {
	return c.doer()
}
