package config

import (
	"encoding/hex"
	"os"
)

type (
	hasSyncKey struct {
		syncKeyName string
		syncKey     []byte
	}
	hasAesKey struct {
		aesKeyName string
		aesKey     []byte
	}
	hasPrivateKey struct {
		privateKey string
	}
	hasPublicKey struct {
		publicKey string
	}
)

func (s *hasSyncKey) loadSyncKey() (err error) {
	if s.syncKeyName == "" {
		return nil
	}
	s.syncKey, err = os.ReadFile(s.syncKeyName)
	return
}

func (s *hasSyncKey) GetSyncKey() []byte {
	return s.syncKey
}

func (a *hasAesKey) decodeAesKey() (err error) {
	if a.aesKeyName == "" {
		return nil
	}
	a.aesKey, err = hex.DecodeString(a.aesKeyName)
	return
}

func (a *hasAesKey) GetAesKey() []byte {
	return a.aesKey
}

func (p *hasPrivateKey) GetPrivateKey() string {
	return p.privateKey
}

func (p *hasPublicKey) GetPublicKey() string {
	return p.publicKey
}
