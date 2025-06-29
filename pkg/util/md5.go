package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// MD5 加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	passwordMdsBys := h.Sum(nil)
	return hex.EncodeToString(passwordMdsBys)
}

// SHA1 加密
func SHA1(str string) string {
	fmt.Println("str:", str)
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

// HMACSHA1 加密
func HMACSHA1(keyStr string, data string) string {
	mac := hmac.New(sha1.New, []byte(keyStr))
	mac.Write([]byte(data))
	srcBytes := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(srcBytes)
}
