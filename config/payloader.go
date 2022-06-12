package config

import (
	"github.com/iv-menshenin/hideme/exec"
)

type payload struct {
	value string
}

func (p *payload) GetPayload() string {
	return p.value
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
