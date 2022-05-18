package crypt

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncryptDataAES(t *testing.T) {
	t.Parallel()
	type args struct {
		data []byte
		key  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "short",
			args: args{
				data: []byte("кому на руси жить хорошо"),
				key:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			},
			wantErr: false,
		},
		{
			name: "fish",
			args: args{
				data: []byte(fishText1),
				key:  bytes.Repeat([]byte{1, 255}, 16),
			},
			wantErr: false,
		},
		{
			name: "bad key",
			args: args{
				data: []byte(fishText1),
				key:  []byte{1, 3},
			},
			wantErr: true,
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var src = make([]byte, len(test.args.data))
			copy(src, test.args.data)
			enc, err := EncryptDataAES(src, test.args.key)
			if (err != nil) != test.wantErr {
				t.Errorf("EncryptDataAES() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if err != nil {
				return
			}
			got, err := DecryptDataAES(enc, test.args.key)
			if (err != nil) != test.wantErr {
				t.Errorf("EncryptDataAES() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.args.data) {
				t.Errorf("matching error\nwant: %s\n got: %s", test.args.data, got)
			}
		})
	}
}
