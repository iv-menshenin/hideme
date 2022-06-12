package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/iv-menshenin/hideme/carrier"
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

func payloadFromQuery(q Query, keyName string) (*payload, error) {
	r, name, err := q.ByteVal(keyName)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	msg, err := message.NewFromBytes(name, data)
	if err != nil {
		return nil, err
	}
	return &payload{
		fileName: name,
		message:  msg,
	}, nil
}

type input struct {
	value string
	carr  carrier.Carrier
}

func (i *input) prepare() (err error) {
	i.carr, err = carrier.NewCarrierFromFile(i.value)
	return
}

func (i *input) GetInput() carrier.Carrier {
	return i.carr
}

func inputFromQuery(q Query, keyName string) (*input, error) {
	r, name, err := q.ByteVal(keyName)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	carr, err := carrier.New(r)
	if err != nil {
		return nil, err
	}
	return &input{
		value: name,
		carr:  carr,
	}, nil
}

type output struct {
	path  string
	files []string
}

func (o *output) SaveFile(fileName string) (io.WriteCloser, error) {
	fileName = o.path + fileName
	f, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot create new file `%s`: %w", fileName, err)
	}
	o.files = append(o.files, fileName)
	return f, nil
}

func (o *output) getFiles() []string {
	return o.files
}
