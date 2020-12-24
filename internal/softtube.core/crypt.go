package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Crypt : Encrypts and decrypts stuff
type Crypt struct {
}

// Not so secret key!
var key []byte = []byte("SoftTeam and SoftTube are good!!") // 32 bytes

// https://stackoverflow.com/questions/18817336/golang-encrypting-a-string-with-aes-and-base64

// Encrypt : Encrypts a string
func (c Crypt) Encrypt(plainText string) (string, error) {
	text := []byte(plainText)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt : Decrypts an encrypted string
func (c Crypt) Decrypt(encryptedString string) (string, error) {
	text, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(text) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return "", err
	}
	return string(data), nil
}