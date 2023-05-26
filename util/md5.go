package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// Md5Encode 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	str := h.Sum(nil)
	return hex.EncodeToString(str)
}

// MD5Encode 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// MakePassword 加密操作
func MakePassword(plainTxt, salt string) string {
	return Md5Encode(plainTxt + salt)
}

// ValidPassword 解密
func ValidPassword(plainTxt, salt string, password string) bool {
	return Md5Encode(plainTxt+salt) == password
}
