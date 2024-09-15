package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// Key should be 16, 24, or 32 bytes long
var key = []byte("your-32-byte-long-secret-key-1234")

// Encrypt encrypts plaintext using AES encryption
func Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES encryption
func Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
