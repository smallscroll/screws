package screws

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

//ITinyTools 小工具接口
type ITinyTools interface {
	DigitalCaptcha() string
	UUIDV4() string
	SHA256OfString(str string) string
	SHA1OfString(str string) string
	SHA256OfFile(fileHeader *multipart.FileHeader) (string, error)
	SHA1OfFile(fileHeader *multipart.FileHeader) (string, error)
	StringsToFloats(strings ...string) ([]float64, error)
	CheckText(str, exp string) bool
	CheckDatetime(str ...string) ([]*time.Time, error)
	CheckTimestamp(str ...string) ([]*time.Time, error)
	CheckDatetimeToTimestamp(str ...string) ([]int64, error)
	CheckUserAgent(userAgent string) string
}

//NewTinyTools 初始化小工具
func NewTinyTools() ITinyTools {
	return &tinyTools{}
}

//tinyTools  小工具
type tinyTools struct {
}

//DigitalCaptcha 六位数字验证码
func (t *tinyTools) DigitalCaptcha() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Intn(899999)+100000)
}

//UUID  UUID/V4
func (t *tinyTools) UUIDV4() string {
	s := strings.Split(uuid.NewV4().String(), "-")
	return strings.Join(s, "")
}

//HashOfString 计算字符串哈希
func (t *tinyTools) SHA256OfString(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}

//HashOfString 计算字符串哈希
func (t *tinyTools) SHA1OfString(str string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
}

//HashOfFile 计算文件哈希
func (t *tinyTools) SHA256OfFile(fileHeader *multipart.FileHeader) (string, error) {
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

//HashOfFile 计算文件哈希
func (t *tinyTools) SHA1OfFile(fileHeader *multipart.FileHeader) (string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	sha := sha1.New()
	_, err = io.Copy(sha, src)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

//StringsToFloats 字符串转浮点数
func (t *tinyTools) StringsToFloats(strings ...string) ([]float64, error) {
	var slice []float64
	for _, v := range strings {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		slice = append(slice, f)
	}
	return slice, nil
}

//CheckText 检查普通字符串格式
func (t *tinyTools) CheckText(str, exp string) bool {
	if !regexp.MustCompile(exp).MatchString(str) {
		return false
	}
	return true
}

//CheckDatetime 检查日期时间字符串格式
func (t *tinyTools) CheckDatetime(str ...string) ([]*time.Time, error) {
	var times []*time.Time
	for _, v := range str {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if err != nil {
			return nil, err
		}
		times = append(times, &t)
	}
	return times, nil
}

//CheckTimestamp 检查时间戳字符串格式
func (t *tinyTools) CheckTimestamp(str ...string) ([]*time.Time, error) {
	var times []*time.Time
	for _, v := range str {
		datetime, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		t := time.Unix(int64(datetime), 0)
		times = append(times, &t)
	}
	return times, nil
}

//CheckDatetimeToTimestamp 检查日期时间字符串格式并转换为时间戳
func (t *tinyTools) CheckDatetimeToTimestamp(str ...string) ([]int64, error) {
	var timestamps []int64
	for _, v := range str {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", v, time.Local)
		if err != nil {
			return nil, err
		}
		timestamps = append(timestamps, t.Unix())
	}
	return timestamps, nil
}

//CheckUserAgent 检查用户客户端类型
func (t *tinyTools) CheckUserAgent(userAgent string) string {
	if regexp.MustCompile(`^(uni-app)+$`).MatchString(userAgent) {
		return "app"
	}
	if regexp.MustCompile(`^(Mozilla)+|(AppleWebKit)+|(Chrome)+|(Safari)+|(Edge)+$`).MatchString(userAgent) {
		if regexp.MustCompile(`^(iPhone)+|(Android)+$`).MatchString(userAgent) {
			return "mobile"
		}
		return "desktop"
	}
	return "other"
}
