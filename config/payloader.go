package config

import (
	"github.com/iv-menshenin/hideme/exec"
	"github.com/iv-menshenin/hideme/message"
)

type payload struct {
	fileName string
	message  *message.Message
}

func (p *payload) prepare() (err error) {
	p.message, err = message.NewFromFile(p.fileName)
	return
}

func (p *payload) GetPayload() *message.Message {
	return p.message
}

type input struct {
	value string
	carr  exec.Carrier
}

func (i *input) prepare() (err error) {
	i.carr, err = exec.NewCarrierFromFile(i.value)
	return
}

func (i *input) GetInput() exec.Carrier {
	return i.carr
}

type output struct {
	value string
}

func (i *output) GetOutput() string {
	return i.value
}
