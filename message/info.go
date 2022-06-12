package message

import (
	"bytes"
	"fmt"
)

func (m *Message) FileName() string {
	var n = 0
	for ; n < fileNameMaxLen; n++ {
		if m.fileName[n] == 0 {
			break
		}
	}
	return string(m.fileName[:n])
}

func (m *Message) Content() []byte {
	return m.content
}

func (m *Message) fillFileName(r *bytes.Reader) error {
	var fileNameSzBy [8]byte
	_, err := r.Read(fileNameSzBy[:])
	if err != nil {
		return fmt.Errorf("cannot extract file name size information: %w", err)
	}

	fileNameSize := b64int(fileNameSzBy)
	if fileNameSize > fileNameMaxLen || fileNameSize < 0 {
		return fmt.Errorf("wrong format")
	}
	if _, err = r.Read(m.fileName[:fileNameSize]); err != nil {
		return fmt.Errorf("cannot extract file name: %w", err)
	}

	return nil
}

func (m *Message) fillFileContent(r *bytes.Reader) error {
	var fileSzBy [8]byte
	if _, err := r.Read(fileSzBy[:]); err != nil {
		return fmt.Errorf("cannot extract file size information: %w", err)
	}

	m.fileSize = b64int(fileSzBy)
	m.content = make([]byte, m.fileSize)
	n, err := r.Read(m.content[:])
	if err != nil {
		return fmt.Errorf("cannot extract file data: %w", err)
	}
	if n != int(m.fileSize) {
		return fmt.Errorf("cannot extract file data: need %d bytes, but read %d", m.fileSize, n)
	}

	return nil
}

func (m *Message) setFileName(name string) {
	nmBytes := []byte(name)
	if len(nmBytes) > fileNameMaxLen {
		copy(m.fileName[:fileNameMaxLen], nmBytes[len(nmBytes)-fileNameMaxLen:])
		return
	}
	copy(m.fileName[:fileNameMaxLen], nmBytes[:])
}

func (m *Message) fileNameSize() int {
	var n = 0
	for ; n < fileNameMaxLen; n++ {
		if m.fileName[n] == 0 {
			break
		}
	}
	return n
}
