package message

type (
	serialized struct {
		data []byte
	}
)

func (m *message) Serialize() []uint8 {
	var s serialized
	s.putChunk(m.fileName[:])
	fs := encode64(m.fileSize)
	s.putChunk(fs[:])
	s.putChunk(m.content)
	return s.getWholeData()
}

func (m *message) Deserialize(data []uint8) {
	var s = serialized{data: data}
	m.setFileName(string(s.pullChunk()))
	fsB := s.pullChunk()
	var fsV [16]byte
	copy(fsV[:], fsB[:])
	m.fileSize = decode64(fsV)
	m.content = s.pullChunk()
}

const (
	rcsSize  = 2
	smtSize  = 4
	smtValue = 15
)

func (s *serialized) getWholeData() []byte {
	return s.data
}

func (s *serialized) putChunk(b []byte) {
	data := encodeB(b)
	i := int64(len(data))
	sz := encode64(i)
	s.data = append(s.data, sz[:]...)
	s.data = append(s.data, data...)
}

func (s *serialized) pullChunk() []byte {
	var i [16]byte
	copy(i[:], s.data[:len(i)])
	l := decode64(i)
	var data = make([]byte, l)
	copy(data[:], s.data[16:len(data)+16])
	s.data = s.data[len(data)+16:]
	return decodeB(data)
}

func encode64(i int64) (result [16]byte) {
	for n := 0; n < 16; n++ {
		m := i & smtValue
		result[n] = byte(m)
		i = i >> smtSize
	}
	return
}

func decode64(d [16]byte) (result int64) {
	for n := 15; n > 0; n-- {
		result = result | int64(d[n])
		result = result << smtSize
	}
	return result | int64(d[0])
}

func encodeB(b []byte) []byte {
	var result = make([]byte, 0, len(b)*rcsSize)
	for n := 0; n < len(b); n++ {
		result = append(result, b[n]&smtValue)
		result = append(result, (b[n]>>smtSize)&smtValue)
	}
	return result
}

func decodeB(b []byte) []byte {
	var result = make([]byte, 0, len(b)/rcsSize)
	for n := 0; n < len(b); n = n + 2 {
		result = append(result, b[n]+b[n+1]<<smtSize)
	}
	return result
}
