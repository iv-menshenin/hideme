package main

import (
	"flag"
	"log"
	"os"

	"github.com/iv-menshenin/hideme/carrier"
	"github.com/iv-menshenin/hideme/message"
)

type cmd struct {
	command string
	input   string
	payload string
	output  string
}

const (
	cmdInject  = "inject"
	cmdExtract = "extract"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("available commands: %s, %s", cmdInject, cmdExtract)
	}

	var config = parseCmd(os.Args[1], os.Args[2:])
	switch config.command {
	case cmdInject:
		if config.input == "" {
			log.Fatal("`carrier` parameter cannot be empty")
		}
		if config.output == "" {
			log.Fatal("`out` parameter cannot be empty")
		}
		if config.payload == "" {
			log.Fatal("`payload` parameter cannot be empty")
		}
		inject(config.payload, config.input, config.output)

	case cmdExtract:
		if config.input == "" {
			log.Fatal("`input` parameter cannot be empty")
		}
		extract(config.input)

	}
}

func parseCmd(toDo string, args []string) cmd {
	var command = cmd{command: toDo}
	switch toDo {
	case cmdInject:
		fs := flag.NewFlagSet(toDo, flag.ExitOnError)
		fs.StringVar(&command.input, "carrier", "", "A PNG file that will carry the valuable information")
		fs.StringVar(&command.payload, "payload", "", "The file you want to hide from prying eyes")
		fs.StringVar(&command.output, "out", "", "The final file, which does not differ from the original file. But it contains encrypted information")
		if err := fs.Parse(args); err != nil {
			log.Fatal(err)
		}

	case cmdExtract:
		fs := flag.NewFlagSet(toDo, flag.ExitOnError)
		fs.StringVar(&command.input, "input", "", "A file that carries hidden information")
		if err := fs.Parse(args); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("available commands: %s, %s\nunknown command: %s", cmdInject, cmdExtract, toDo)
	}
	return command
}

func inject(payload, source, outFile string) {
	msg, err := message.New(payload)
	if err != nil {
		log.Fatalf("cannot prepare msg: %s", err)
	}

	carr, err := getCarrier(source)
	if err != nil {
		log.Fatalf("cannot prepare carrier file: %s", err)
	}
	secretData := msg.Serialize()
	if err = carr.Inject(secretData); err != nil {
		log.Fatalf("cannot inject secret data to image: %s", err)
	}
	if err = carr.SaveTo(outFile); err != nil {
		log.Fatalf("cannot save image file: %s", err)
	}
}

func extract(payload string) {
	carr, err := getCarrier(payload)
	if err != nil {
		log.Fatalf("cannot prepare carrier file: %s", err)
	}
	data := carr.GetPayload()
	msg := message.FromData(data)
	f, err := os.Create("./_" + msg.FileName())
	if err != nil {
		log.Fatalf("cannot create new file `%s`: %s", msg.FileName(), err)
	}
	defer f.Close()
	_, err = f.Write(msg.Content())
	if err != nil {
		log.Fatalf("cannot write data to file `%s`: %s", msg.FileName(), err)
	}
}

type injector interface {
	Inject([]uint8) error
	GetPayload() []uint8
	SaveTo(string) error
}

func getCarrier(fileName string) (injector, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return carrier.New(f)
}
