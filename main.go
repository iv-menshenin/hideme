package main

import (
	"fmt"
	"github.com/iv-menshenin/hideme/http"
	"log"
	"os"
	"strings"

	"github.com/iv-menshenin/hideme/config"
)

const (
	cmdInject   = "inject"
	cmdExtract  = "extract"
	cmdGenerate = "keys"
	cmdServer   = "server"
)

var available = []string{cmdInject, cmdExtract, cmdGenerate, cmdServer}

type (
	Configurator interface {
		Parse() error
		Validate() error
		Execute() error
	}
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("available commands: %s", strings.Join(available, ", "))
	}

	cmd, err := parseCmd(os.Args[1], os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
	if err = cmd.Validate(); err != nil {
		log.Fatal(err)
	}
	if err = cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func parseCmd(toDo string, args []string) (Configurator, error) {
	var cmd Configurator

	switch toDo {
	case cmdInject:
		cmd = config.NewInjector(args)

	case cmdExtract:
		cmd = config.NewExtractor(args)

	case cmdGenerate:
		cmd = config.NewGenerator(args)

	case cmdServer:
		cmd = config.NewServer(args, http.Handler)

	case "--help", "-h", "help":
		fmt.Print(helpInformation)
		os.Exit(0)

	default:
		return nil, fmt.Errorf("available commands: %s\nunknown command: %s", strings.Join(available, ", "), toDo)
	}

	err := cmd.Parse()
	if err != nil {
		return nil, fmt.Errorf("can't parse parameters: %s", err)
	}
	return cmd, nil
}
