package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func SignData(data []byte, privateKey string) ([]byte, error) {
	private, err := getPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot parse private key: %w", err)
	}
	sign, err := rsa.SignPSS(rand.Reader, private, signHashFn, hashData(data), nil)
	if err != nil {
		return nil, fmt.Errorf("error while signing: %w", err)
	}
	return sign, nil
}

func SignVerify(data, sign []byte, publicKey string) error {
	public, err := getPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("cannot parse public key `%s`: %w", publicKey, err)
	}
	err = rsa.VerifyPSS(public, signHashFn, hashData(data), sign, nil)
	if err != nil {
		return fmt.Errorf("error while sign checking: %w", err)
	}
	return nil
}
