package math

import (
	"math/rand"
	"time"
)

// rand seeds
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// Random16Str 返回指定长度的随机16进制数字字符串
func Random16Str(length int) string {
	return RandomStr(length, 16)
}

// Random62Str 返回指定长度的随机62进制数字字符串
func Random62Str(length int) string {
	return RandomStr(length, 62)
}

// RandomStr 返回指定长度和进制的随机数字字符串
func RandomStr(length, base int) string {
	if length <= 0 || base <= 0 || base > len(BaseNumber) {
		return ""
	}
	str := make([]byte, length)
	for i := 0; i < length; i++ {
		str[i] = BaseNumber[seed.Intn(base)]
	}
	return string(str)
}
