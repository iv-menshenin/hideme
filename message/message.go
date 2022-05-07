package message

import (
	"io/ioutil"
	"os"
	"path"
)

type (
	message struct {
		fileName fileName
		fileSize int64
		content  []byte
	}
	fileName [fileNameMaxLen]byte
)

const fileNameMaxLen = 64

func New(fileName string) (*message, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	m := message{
		fileSize: stat.Size(),
	}
	m.content, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	m.setFileName(path.Base(fileName))
	return &m, nil
}

func FromData(data []byte) *message {
	var m message
	m.Deserialize(data)
	return &m
}
