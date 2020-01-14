package screws

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

//ICrypto 加密器接口
type ICrypto interface {
	GetKey() string
	AESCTREncrypt(plaintextString string) (string, error)
	AESCTRDecrypt(ciphertextString string) (string, error)
}

//NewCrypto 初始化加密器
func NewCrypto(key string) ICrypto {
	return &crypto{
		Key: key,
	}
}

//crypto 加密器
type crypto struct {
	Key string
}

//GetKey 查询密钥
func (c *crypto) GetKey() string {
	return c.Key
}

//AESCTREncrypt AES-CTR加密
func (c *crypto) AESCTREncrypt(plaintextString string) (string, error) {
	plaintext := plaintextString
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
func (c *crypto) AESCTRDecrypt(ciphertextString string) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextString)
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
