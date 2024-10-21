package util

import (
	"math/rand"
	"time"
)

func Random(count int, availableChars string) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune(availableChars)
	charsLen := len(chars)
	var result []rune
	for i := 0; i < count; i++ {
		result = append(result, chars[random.Intn(charsLen)])
	}
	return string(result)
}

//生成随机字符串,包含小写字母
func RandomStr(count int) string {
	return Random(count, "abcdefghijklmnopqrstuvwxyz")
}

//生成随机数字字符串
func RandomNumStr(count int) string {
	return Random(count, "0123456789")
}

//生成随机字符串,包含在ascii 33- 126的字符
func RandomAscii(count int) string {
	start := 33
	end := 127
	bytes := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := end - start
	for len(bytes) < count {
		num := r.Intn(interval) + start
		bytes = append(bytes, byte(num))
	}
	return string(bytes)
}

//随机生成count个数字,数字大小范围:start - end(不包含)
func RandomNumbers(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	var nums []int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := end - start
	for len(nums) < count {
		num := r.Intn(interval) + start
		nums = append(nums, num)
	}

	return nums
}
