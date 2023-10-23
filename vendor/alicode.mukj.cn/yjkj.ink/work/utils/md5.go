package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//返回一个32位md5加密后的字符串
func MD5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

//返回一个16位md5加密后的字符串
func Get16MD5(data string) string {
	return MD5(data)[8:24]
}

func test() {
	
}