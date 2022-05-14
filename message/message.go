package message

import (
	"bytes"
	"io/ioutil"
	"math"
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

func (m *message) Encode() []byte {
	var result bytes.Buffer

	fileNameSize := m.fileNameSize()
	fileNameSzBy := int64b(int64(fileNameSize))
	result.Write(fileNameSzBy[:])
	result.Write(m.fileName[:fileNameSize])

	fileSzBy := int64b(m.fileSize)
	result.Write(fileSzBy[:])
	result.Write(m.content)

	return result.Bytes()
}

func Decode(data []byte) (*message, error) {
	var m message
	var r = bytes.NewReader(data)

	if err := m.fillFileName(r); err != nil {
		return nil, err
	}

	if err := m.fillFileContent(r); err != nil {
		return nil, err
	}

	return &m, nil
}

func int64b(i int64) (result [8]byte) {
	for n := 0; n < len(result); n++ {
		m := i & math.MaxUint8
		result[n] = byte(m)
		i = i >> 8
	}
	return
}

func b64int(d [8]byte) (result int64) {
	for n := 7; n > 0; n-- {
		result = result | int64(d[n])
		result = result << 8
	}
	return result | int64(d[0])
}
