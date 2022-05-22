package crypt

import (
	"bytes"
	"crypto/aes"
	"strconv"
	"strings"
)

func EncryptDataAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var chainSize = block.BlockSize()
	var l = len(data)
	var sl = strconv.FormatInt(int64(l), 10)
	sl = strings.Repeat("0", chainSize-len(sl)) + sl

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

func DecryptDataAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var bl = block.BlockSize()
	var lz = make([]byte, bl)
	block.Decrypt(lz, data[:bl])

	dataLen, err := strconv.ParseInt(string(lz), 10, 64)
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
