package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"strconv"
)

func EncryptDataAES(data []byte, key []byte) ([]byte, error) {
	aesEncoder, err := newAES(key)
	if err != nil {
		return nil, err
	}
	chainSize := aesEncoder.blockSize()
	infoBlock := newSizeInfoChunk(len(data), chainSize)
	data = alignDataBy(data, chainSize)
	encrypted := make([]byte, len(infoBlock)+len(data))
	if err = aesEncoder.encode(encrypted[0:len(infoBlock)], infoBlock); err != nil {
		return nil, err
	}

	for n := 0; n < len(data)/chainSize; n++ {
		var dst, src = encrypted[(n+1)*chainSize : (n+2)*chainSize], data[n*chainSize : (n+1)*chainSize]
		if err = aesEncoder.encode(dst, src); err != nil {
			return nil, err
		}
	}
	return encrypted, nil
}

func DecryptDataAES(data []byte, key []byte) ([]byte, error) {
	aesEncoder, err := newAES(key)
	if err != nil {
		return nil, err
	}
	chainSize := aesEncoder.blockSize()
	decrypted := make([]byte, len(data))
	for n := 0; n < len(data)/chainSize; n++ {
		var dst, src = decrypted[n*chainSize : (n+1)*chainSize], data[n*chainSize : (n+1)*chainSize]
		if err = aesEncoder.decode(dst, src); err != nil {
			return nil, err
		}
	}

	dataLen, err := getDataSize(decrypted, chainSize)
	if err != nil {
		return nil, err
	}
	return decrypted[chainSize : chainSize+dataLen], nil
}

func newSizeInfoChunk(dataLen int, chainSize int) (result []byte) {
	result = make([]byte, chainSize)
	sl := []byte(
		strconv.FormatInt(int64(dataLen), 10),
	)
	noiseLen := chainSize - len(sl)

	for i := 0; i < noiseLen; i++ {
		result[i] = getRandByteNaN()
	}
	copy(result[noiseLen:], sl)
	return
}

func getRandByteNaN() byte {
	var b = make([]byte, 1)
	for {
		if _, err := rand.Read(b[:1]); err != nil {
			panic(err) // newer happens
		}
		if b[0] < '0' || b[0] > '9' {
			return b[0]
		}
	}
}

func alignDataBy(data []byte, alignment int) []byte {
	var resultChains = len(data) / alignment
	if resultChains*alignment < len(data) {
		addSz := len(data) - resultChains*alignment
		// add some salt to last chain
		if len(data) > alignment {
			data = append(data, data[len(data)-alignment-addSz:len(data)-addSz]...)
		} else {
			data = append(data, bytes.Repeat([]byte{data[0]}, addSz)...)
		}
	}
	return data
}

func getDataSize(data []byte, chainSize int) (int, error) {
	var lenStart = 0
	for ; lenStart < chainSize; lenStart++ {
		if data[lenStart] >= '0' && data[lenStart] <= '9' {
			break
		}
	}
	dataLen, err := strconv.ParseInt(string(data[lenStart:chainSize]), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(dataLen), nil
}

type encoder struct {
	cipher cipher.Block
	initVc []byte
}

func newAES(key []byte) (*encoder, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	enc := encoder{
		cipher: block,
		initVc: make([]byte, block.BlockSize()),
	}
	return &enc, nil
}

func (e *encoder) blockSize() int {
	return e.cipher.BlockSize()
}

func (e *encoder) encode(dst, src []byte) (err error) {
	c := newChain(src, dst, e.initVc)
	if err = c.mixInput(); err != nil {
		return
	}
	return c.encrypt(e.cipher)
}

func (e *encoder) decode(dst, src []byte) (err error) {
	c := newChain(src, dst, e.initVc)
	if err = c.decrypt(e.cipher); err != nil {
		return
	}
	return c.mixOutput()
}

type chain struct {
	initV []byte
	inp   []byte
	out   []byte
	inter []byte
}

func newChain(inp, out, initV []byte) chain {
	return chain{
		initV: initV,
		inp:   inp,
		out:   out,
		inter: make([]byte, len(inp)),
	}
}

func (c *chain) mixInput() error {
	return xorData(c.inp, c.initV, c.inter)
}

func (c *chain) encrypt(crp interface{ Encrypt(dst, src []byte) }) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("encryption error: %v", r)
		}
	}()
	crp.Encrypt(c.out, c.inter)
	return xorData(c.inp, c.out, c.initV)
}

func (c *chain) decrypt(crp interface{ Decrypt(dst, src []byte) }) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("decryption error: %v", r)
		}
	}()
	crp.Decrypt(c.inter, c.inp)
	return xorData(c.inter, c.initV, c.out)
}

func (c *chain) mixOutput() error {
	return xorData(c.inp, c.out, c.initV)
}

func xorData(a, b, c []byte) error {
	if len(a) != len(b) || len(b) != len(c) {
		return fmt.Errorf("must be same len, but got %d, %d and %d", len(a), len(b), len(c))
	}
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i]
	}
	return nil
}
