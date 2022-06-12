package config

import (
	"flag"
	"fmt"
)

type (
	server struct {
		httpPort
	}
	httpPort struct {
		port int
	}
)

func (s *server) initCmdParameters() parser {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	fs.IntVar(&s.port, "port", 8080, "The port that serves requests.")
	return fs
}

func (s *server) validate() error {
	if s.port < 80 || s.port > 65535 {
		return fmt.Errorf("port value expects from 80 tot 65535, got: %d", s.port)
	}
	return nil
}

func (s *httpPort) GetPort() int {
	return s.port
}

func (s *server) files() []string {
	return nil
}
