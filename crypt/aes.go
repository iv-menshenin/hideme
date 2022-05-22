package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"strconv"
)

func EncryptDataAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var chainSize = block.BlockSize()
	var l = len(data)
	var sl = []byte(strconv.FormatInt(int64(l), 10))
	for i := len(sl); i < chainSize; i++ {
		if data[i] >= '0' || data[i] <= '9' {
			sl = append([]byte{getRandByteNaN()}, sl...)
			continue
		}
		sl = append([]byte{data[i]}, sl...)
	}

	var result = make([]byte, len(data)+chainSize*2)
	var resultChains = len(data) / chainSize
	if resultChains*chainSize < len(data) {
		addSz := len(data) - resultChains*chainSize
		resultChains++
		// add some salt to last chain
		if len(data) > chainSize {
			data = append(data, data[len(data)-chainSize-addSz:len(data)-addSz]...)
		} else {
			data = append(data, bytes.Repeat([]byte{data[0]}, addSz)...)
		}
	}
	var resultLen = (resultChains + 1) * chainSize
	result = result[:resultLen]
	block.Encrypt(result[:chainSize], []byte(sl))
	var n = 1
	for {
		block.Encrypt(result[n*chainSize:(n+1)*chainSize], data[(n-1)*chainSize:n*chainSize])
		n++
		if n > resultChains {
			break
		}
	}

	return result, nil
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

func DecryptDataAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var bl = block.BlockSize()
	var lz = make([]byte, bl)
	block.Decrypt(lz, data[:bl])
	var lenStart = 0
	for ; lenStart < len(lz); lenStart++ {
		if lz[lenStart] >= '0' && lz[lenStart] <= '9' {
			break
		}
	}

	dataLen, err := strconv.ParseInt(string(lz[lenStart:]), 10, 64)
	if err != nil {
		return nil, err
	}
	var result = make([]byte, len(data))
	var resultChains = len(data) / bl

	var n = 0
	for {
		block.Decrypt(result[n*bl:(n+1)*bl], data[(n+1)*bl:(n+2)*bl])
		n++
		if n > resultChains-2 {
			break
		}
	}
	return result[:dataLen], nil
}
