package util

import (
	"crypto/sha256"
	"fmt"
)

//SHA256摘要计算,返回计算后的16进制字符串
func SHA256Hex(plainStr string) string {
	sum := sha256.Sum256([]byte(plainStr))
	return fmt.Sprintf("%x", sum)
}
