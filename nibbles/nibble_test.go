package nibbles

import (
	"reflect"
	"testing"
)

func TestNewNibbles3(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name   string
		input  []byte
		output []byte
	}
	var testCases = []testCase{
		{
			name:   "one byte",
			input:  []byte{255},
			output: []byte{7, 7, 3},
		},
		{
			name:   "255 fill",
			input:  []byte{255, 255, 255},
			output: []byte{7, 7, 7, 7, 7, 7, 7, 7},
		},
		{
			name:   "chess",
			input:  []byte{170, 170, 170, 170, 170, 170},
			output: []byte{2, 5, 2, 5, 2, 5, 2, 5, 2, 5, 2, 5, 2, 5, 2, 5},
		},
		{
			name:   "char",
			input:  []byte("f"),
			output: []byte{0b110, 0b100, 0b01},
		},
		{
			name:   "expression",
			input:  []byte("foo"),
			output: []byte{0b110, 0b100, 0b101, 0b111, 0b110, 0b110, 0b011, 0b011},
		},
	}
	for i := range testCases {
		test := testCases[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var got []byte
			n := New(3, test.input)
			for {
				b, ok := n.Next()
				if !ok {
					break
				}
				got = append(got, b)
			}
			if !reflect.DeepEqual(test.output, got) {
				t.Errorf("matching error\nwant: %v\n got: %v", test.output, got)
			}
		})
	}
}

func TestNibblingAndConvert(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name  string
		size  int
		input []byte
	}
	var testCases = []testCase{
		{
			name:  "one byte 3",
			size:  3,
			input: []byte{255},
		},
		{
			name:  "one byte 4",
			size:  4,
			input: []byte{255},
		},
		{
			name:  "255 fill 3",
			size:  3,
			input: []byte{255, 255, 255},
		},
		{
			name:  "255 fill 6",
			size:  6,
			input: []byte{255, 255, 255},
		},
		{
			name:  "chess 3",
			size:  3,
			input: []byte{170, 170, 170, 170, 170, 170},
		},
		{
			name:  "chess 5",
			size:  5,
			input: []byte{170, 170, 170, 170, 170, 170},
		},
		{
			name:  "expression 3",
			size:  3,
			input: []byte("foo, bar. let`s play that game"),
		},
		{
			name:  "expression 5",
			size:  5,
			input: []byte("foo, bar. let`s play that game"),
		},
	}
	for i := range testCases {
		test := testCases[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var got []byte
			n := New(test.size, test.input)
			for {
				b, ok := n.Next()
				if !ok {
					break
				}
				got = append(got, b)
			}
			got = Convert(got, test.size)
			if !reflect.DeepEqual(test.input, got) {
				t.Errorf("matching error\nwant: %v\n got: %v", test.input, got)
			}
		})
	}
}
