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

func (s *hasSyncKey) loadSyncKey(syncKey string) (err error) {
	s.syncKey, err = os.ReadFile(syncKey)
	return
}

func (s *hasSyncKey) GetSyncKey() []byte {
	return s.syncKey
}

func (a *hasAesKey) decodeAesKey(aesKey string) (err error) {
	a.aesKey, err = hex.DecodeString(aesKey)
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
