package crypt

import (
	"crypto"
	"crypto/sha512"
)

const signHashFn = crypto.SHA512

func hashData(data []byte) []byte {
	hash := sha512.Sum512(data)
	return hash[:]
}
