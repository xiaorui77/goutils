package math

import (
	"math/rand"
	"time"
)

const number = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_+=!@#$%^&*()"

// rand seeds
var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// Random16Str 返回指定长度的随机16进制数字字符串
func Random16Str(len int) string {
	return RandomStr(len, 16)
}

// Random62Str 返回指定长度的随机62进制数字字符串
func Random62Str(len int) string {
	return RandomStr(len, 62)
}

// RandomStr 返回指定长度和进制的随机数字字符串
func RandomStr(len, base int) string {
	if len <= 0 || base <= 0 || base > 76 {
		return ""
	}
	str := make([]byte, len)
	for i := 0; i < len; i++ {
		str[i] = number[seed.Intn(base)]
	}
	return string(str)
}
