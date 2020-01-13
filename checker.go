package screws

import (
	"regexp"
	"strconv"
	"time"
)

var regexpItems = map[string]string{
	"email":          `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`, //邮箱格式
	"mobile":         `^1[3-9][0-9]{9}$`,                              //手机格式
	"password":       `^([!-~]{8,30})$`,                               //密码格式：8至30位字母/数字/特殊符号
	"bankCard":       `^[0-9]{12,30}$`,                                //银行卡格式：12至30位数字
	"alias":          `^[a-zA-Z0-9\p{Han}]{2,18}$`,                    //用户名格式：2至16位英文或汉字或数字
	"identityName":   `^[a-zA-Z\p{Han}]{2,18}$`,                       //姓名格式：2至16位英文或汉字
	"identityNumber": `^[0-9a-zA-Z]{8,20}$`,                           //证件号码格式：8至20位数字或字母
}

//CheckText 检查字符串格式
func CheckText(item, str, customRegexp string) bool {
	if customRegexp != "" {
		if !regexp.MustCompile(customRegexp).MatchString(str) {
			return false
		}
	} else {
		v, ok := regexpItems[item]
		if !ok {
			return false
		}
		if !regexp.MustCompile(v).MatchString(str) {
			return false
		}
	}
	return true
}

//CheckTimeString 检查时间字符串格式
func CheckTimeString(str ...string) ([]*time.Time, error) {
	var times []*time.Time
	for _, v := range str {
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return nil, err
		}
		times = append(times, &t)
	}
	return times, nil
}

//CheckTimestamp 检查时间戳格式
func CheckTimestamp(str ...string) ([]*time.Time, error) {
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

//CheckUserAgentType 检查用户客户端类型
func CheckUserAgentType(userAgent string) string {
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
