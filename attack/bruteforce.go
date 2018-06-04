package attack

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
)

func readCiphertext(source string) ([]byte, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ciphertext := make([]byte, 16)
	_, err = io.ReadFull(file, ciphertext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

func writeKey(source string, key []byte) error {
	dir := filepath.Dir(source)
	filename := filepath.Join(dir, "key.bin")

	log.Printf("Writing the key to '%s'", filename)

	err := ioutil.WriteFile(filename, key, 0644)

	return err
}

func Bruteforce(source string) error {
	ciphertext, err := readCiphertext(source)
	if err != nil {
		return err
	}

	expected := []byte("SQLite format 3\x00")
	plaintext := make([]byte, 16)

	for iv := math.MinInt32; iv <= math.MaxInt32; iv++ {
		key := deriveKey(iv)

		cipher, err := aes.NewCipher(key)
		if err != nil {
			return err
		}
		cipher.Decrypt(plaintext, ciphertext)

		if bytes.Equal(plaintext, expected) {
			log.Printf("IV recovered: %d\n", iv)

			err = writeKey(source, key)
			if err != nil {
				return err
			}

			dest := fmt.Sprintf("%s.sqlite", source)
			err = decrypt(source, dest, key)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}
