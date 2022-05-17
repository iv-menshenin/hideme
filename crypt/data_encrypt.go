package crypt

import (
	"crypto/md5"
	"errors"
)

var ErrKeyShortedThanData = errors.New("the length of key must be great or equal than length of data")

func EncryptDecryptData(data []byte, key []byte) error {
	key = moreStrongKey(key)
	if len(key) < len(data) {
		return ErrKeyShortedThanData
	}
	for i, d := range data {
		data[i] = d ^ key[i]
	}
	return nil
}

func moreStrongKey(key []byte) []byte {
	const (
		salt   = 170
		bufLen = 16
	)
	var (
		buf [bufLen * 2]byte
		unf int
		out []byte
	)
	flush := func() {
		unf = 0
		h := md5.Sum(buf[:])
		out = append(out, h[:]...)
	}
	for i, b := range key {
		r := key[len(key)-i-1]
		p := i % bufLen
		buf[p*2] = b
		buf[p*2+1] = b ^ r ^ salt
		unf++
		if (i+1)%bufLen == 0 {
			flush()
		}
	}
	if unf > 0 {
		flush()
	}
	return out
}
