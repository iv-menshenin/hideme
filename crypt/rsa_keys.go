package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
)

const keyBitsSize = 2048

func GenerateKeys() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keyBitsSize)
}

func SaveKeysToFile(private *rsa.PrivateKey, pub, prv io.Writer) error {
	if err := pem.Encode(pub, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&private.PublicKey)}); err != nil {
		return fmt.Errorf("cannot encode public file: %w", err)
	}
	if err := pem.Encode(prv, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(private)}); err != nil {
		return fmt.Errorf("cannot encode private file: %w", err)
	}
	return nil
}

func getPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(bytes)
	if p.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected pem type: %s, expect `RSA PRIVATE KEY`", p.Type)
	}
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func getPublicKey(fileName string) (*rsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	p, _ := pem.Decode(bytes)
	if p.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected pem type: %s, expect `RSA PUBLIC KEY`", p.Type)
	}
	return x509.ParsePKCS1PublicKey(p.Bytes)
}
