package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func Btoa(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Atob(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func Encrypt(key []byte, data []byte) ([]byte, error) {
	// Create a new AES cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a new byte slice to hold the IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	// Fill the IV with random bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Create a new AES cipher block mode of operation using the block and IV
	stream := cipher.NewCFBEncrypter(block, iv)

	// Create a byte slice to hold the ciphertext
	ciphertext := make([]byte, len(data))
	// Encrypt the data
	stream.XORKeyStream(ciphertext, data)

	// Prepend the IV to the ciphertext
	ciphertext = append(iv, ciphertext...)

	return ciphertext, nil
}

func Decrypt(key []byte, data []byte) ([]byte, error) {
	// Create a new AES cipher block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Extract the IV from the ciphertext
	iv := data[:aes.BlockSize]
	// Extract the actual ciphertext
	ciphertext := data[aes.BlockSize:]

	// Create a new AES cipher block mode of operation using the block and IV
	stream := cipher.NewCFBDecrypter(block, iv)

	// Create a byte slice to hold the plaintext
	plaintext := make([]byte, len(ciphertext))
	// Decrypt the data
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

func NewPassword() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create new password"), err)
	}
	return key, nil
}
