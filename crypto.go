package main

import (
	"crypto/aes"
	"io/ioutil"
)

func hashCode(str string) int {
	var hash int32
	for i := 0; i < len(str); i++ {
		hash = 31*hash + int32(str[i])
	}

	return int(hash)
}

func generateKey(iv int) []byte {
	const length = 16
	var pad [length]int8
	var key [length]byte

	pad[0] = int8(iv)
	pad[1] = pad[0] - 71
	pad[2] = pad[1] - 71
	for i := 3; i < length; i++ {
		pad[i] = int8(int(pad[i-3]) ^ int(pad[i-2]) ^ 0xb9 ^ i)
	}

	factor := iv
	if iv > -2 && iv < 2 {
		factor = -313187 + 13819823*iv
	}

	term := -7
	for i := 1; i < length+1; i++ {
		index := i & (length - 1)
		value := int(pad[index])*factor + term

		term = int(int8(value >> 32))
		value = int(int32(value + term))
		if value < term {
			term++
			value++
		}

		value = -value - 2
		pad[index] = int8(value)
	}

	for i := 0; i < length; i++ {
		key[i] = byte(pad[i])
	}

	return key[:]
}

func decrypt(source, dest, ourID, theirID string) error {
	iv := hashCode(ourID + theirID)
	key := generateKey(iv)

	cipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	blockSize := cipher.BlockSize()

	ciphertext, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	length := len(ciphertext)

	plaintext := make([]byte, length)

	for i := 0; i < length; i += blockSize {
		cipher.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
	}

	err = ioutil.WriteFile(dest, plaintext, 0644)

	return err
}
