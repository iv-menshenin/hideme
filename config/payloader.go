package config

type payload struct {
	value string
}

func (p *payload) GetPayload() string {
	return p.value
}

type input struct {
	value string
}

func (i *input) GetInput() string {
	return i.value
}

type output struct {
	value string
}

func (i *output) GetOutput() string {
	return i.value
}
