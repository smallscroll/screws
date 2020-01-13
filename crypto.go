package screws

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

//Crypto ...
type Crypto struct {
	Key        string
	Plaintext  string
	Ciphertext string
}

//AESCTREncrypt AES-CTR加密
func (c *Crypto) AESCTREncrypt() (string, error) {
	plaintext := c.Plaintext
	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len([]byte(plaintext)))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	return hex.EncodeToString(ciphertext), nil
}

//AESCTRDecrypt AES-CTR解密
func (c *Crypto) AESCTRDecrypt() (string, error) {
	ciphertext, err := hex.DecodeString(c.Ciphertext)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	stream := cipher.NewCTR(block, iv)
	plaintext := ciphertext[aes.BlockSize:]
	stream.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])
	return string(plaintext), nil
}
