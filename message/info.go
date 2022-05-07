package message

func (m *message) FileName() string {
	var n = 0
	for ; n < fileNameMaxLen; n++ {
		if m.fileName[n] == 0 {
			break
		}
	}
	return string(m.fileName[:n])
}

func (m *message) Content() []byte {
	return m.content
}

func (m *message) setFileName(name string) {
	nmBytes := []byte(name)
	if len(nmBytes) > fileNameMaxLen {
		copy(m.fileName[:fileNameMaxLen], nmBytes[len(nmBytes)-fileNameMaxLen:])
		return
	}
	copy(m.fileName[:fileNameMaxLen], nmBytes[:])
}
