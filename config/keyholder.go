package config

import (
	"encoding/hex"
	"io/ioutil"
	"os"
)

type (
	hasSyncKey struct {
		syncKeyName string
		syncKey     []byte
	}
	hasAesKey struct {
		keyStr string
		aesKey []byte
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

func syncKeyFromQuery(q Query, keyName string) (*hasSyncKey, error) {
	r, name, err := q.ByteVal(keyName)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &hasSyncKey{
		syncKeyName: name,
		syncKey:     data,
	}, nil
}

func (a *hasAesKey) decodeAesKey() (err error) {
	if a.keyStr == "" {
		return nil
	}
	a.aesKey, err = hex.DecodeString(a.keyStr)
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
