package screws

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"time"
)

//DigitalCaptcha 六位数字验证码
func DigitalCaptcha() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(899999)+100000)
}

//HashOfString 计算字符串哈希
func HashOfString(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

//HashOfFile 计算文件哈希
func HashOfFile(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	sha := sha256.New()
	_, err = io.Copy(sha, src)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}
