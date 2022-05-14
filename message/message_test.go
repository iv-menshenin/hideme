package message

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	const (
		fileName    = "2022_04_29-test-long-filename-09876543221-no-buff-foo-bar-1234567890.txt"
		allowedName = "29-test-long-filename-09876543221-no-buff-foo-bar-1234567890.txt"
	)
	m, err := New("./test/" + fileName)
	if err != nil {
		t.Errorf("cannot create message: %s", err)
	}
	if m.FileName() != allowedName {
		t.Errorf("want: %s, got: %s", allowedName, m.FileName())
	}
	data := m.Encode()
	newM, err := Decode(data)
	if err != nil {
		t.Errorf("decoding error: %s", err)
		return
	}
	if !reflect.DeepEqual(m, newM) {
		t.Errorf("matching error\nwant: %+v\n got: %+v", m, newM)
	}
}
