package message

import (
	"math"
	"reflect"
	"testing"
)

func Test_encode64(t *testing.T) {
	t.Parallel()
	var testCases = []int64{
		0,
		5577006791947779410,
		8674665223082153551,
		6129484611666145821,
		4037200794235010051,
		3916589616287113937,
		6334824724549167320,
		605394647632969758,
		1443635317331776148,
		894385949183117216,
		2775422040480279449,
		math.MaxInt64,
	}
	for _, i := range testCases {
		e := encode64(i)
		d := decode64(e)
		if i != d {
			t.Errorf("matching error: want: %d, got: %d", i, d)
		}
	}
}

func Test_encodeB(t *testing.T) {
	t.Parallel()
	testCases := [][]byte{
		{255, 0, 255, 0, 255, 0, 255},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		[]byte("кому на руси жить хорошо"),
		{},
	}
	for _, i := range testCases {
		e := encodeB(i)
		d := decodeB(e)
		if !reflect.DeepEqual(i, d) {
			t.Errorf("matching error:\nwant: %v\n got: %v", i, d)
		}
	}
}

func Test_serialized_getChunk(t *testing.T) {
	t.Parallel()
	testCases := [][]byte{
		{255, 0, 255, 0, 255, 0, 255},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		[]byte("кому на руси жить хорошо"),
		{},
	}
	var s serialized
	for _, i := range testCases {
		s.putChunk(i)
		d := s.pullChunk()
		if !reflect.DeepEqual(i, d) {
			t.Errorf("matching error:\nwant: %v\n got: %v", i, d)
		}
	}
}
