package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

const keyBitsSize = 2048

func GenerateKeys() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, keyBitsSize)
}

func SaveKeysToFile(private *rsa.PrivateKey, fileName string) error {
	pub, err := os.Create(fmt.Sprintf("%s.pub", fileName))
	if err != nil {
		return fmt.Errorf("cannot create file `%s.pub`: %w", fileName, err)
	}
	defer pub.Close()
	prv, err := os.Create(fmt.Sprintf("%s", fileName))
	if err != nil {
		return fmt.Errorf("cannot create file `%s`: %w", fileName, err)
	}
	defer prv.Close()
	if err = pem.Encode(pub, &pem.Block{Type: "RSA PUBLIC KEY", Bytes: x509.MarshalPKCS1PublicKey(&private.PublicKey)}); err != nil {
		return fmt.Errorf("cannot encode public file: %w", err)
	}
	if err = pem.Encode(prv, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(private)}); err != nil {
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
