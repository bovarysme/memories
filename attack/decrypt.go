package attack

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type decryptor struct {
	r      io.Reader
	w      io.Writer
	cipher cipher.Block
}

func (d *decryptor) decryptChunk(length int) error {
	ciphertext := make([]byte, length)
	_, err := io.ReadFull(d.r, ciphertext)
	if err != nil {
		return err
	}

	blockSize := d.cipher.BlockSize()

	plaintext := make([]byte, length)
	for i := 0; i < length; i += blockSize {
		d.cipher.Decrypt(plaintext[i:i+blockSize], ciphertext[i:i+blockSize])
	}

	// XXX: handle possible panics
	padding := int(plaintext[length-1])
	plaintext = plaintext[:length-padding]

	_, err = d.w.Write(plaintext)

	return err
}

func readChunkLengths(source string) ([]int, error) {
	filename := fmt.Sprintf("%s.extra", source)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	strs := strings.Split(string(data), ",")
	length := len(strs)

	lengths := make([]int, length)
	for i := 0; i < length; i++ {
		lengths[i], err = strconv.Atoi(strs[i])
		if err != nil {
			return nil, err
		}
	}

	return lengths, nil
}

func decrypt(source, dest string, key []byte) error {
	log.Printf("Decrypting '%s' to '%s'\n", source, dest)

	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer writer.Close()

	cipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	d := &decryptor{
		r:      reader,
		w:      writer,
		cipher: cipher,
	}

	lengths, err := readChunkLengths(source)
	if err != nil {
		return err
	}

	for _, length := range lengths {
		err = d.decryptChunk(length)
		if err != nil {
			return err
		}
	}

	return nil
}

func Decrypt(source, dest, ourID, theirID string) error {
	iv := hashCode(ourID + theirID)
	key := deriveKey(iv)

	err := decrypt(source, dest, key)

	return err
}
