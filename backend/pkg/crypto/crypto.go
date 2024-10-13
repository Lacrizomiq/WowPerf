package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	key     []byte
	keyOnce sync.Once
)

func initKey() error {
	var err error
	keyOnce.Do(func() {
		encodedKey := os.Getenv("ENCRYPTION_KEY")
		if encodedKey == "" {
			err = errors.New("ENCRYPTION_KEY environment variable is not set")
			return
		}

		key, err = base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			return
		}

		if len(key) != 32 {
			err = errors.New("ENCRYPTION_KEY must be 32 bytes long (256 bits)")
			return
		}
	})
	return err
}

func Encrypt(plaintext []byte) ([]byte, error) {
	if err := initKey(); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func Decrypt(ciphertext []byte) ([]byte, error) {
	if err := initKey(); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
