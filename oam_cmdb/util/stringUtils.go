package util

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/exp/constraints"
)

const (
	UNSAFE_SQL_CHAR = "`><!%|=*#"
)

//以分隔符拼接整形切片
func JoinInteger[T constraints.Integer](sep string, elems ...T) string {
	var strBuilder strings.Builder
	var lastIndex = len(elems) - 1
	for i, value := range elems {
		strBuilder.WriteString(fmt.Sprintf("%d", value))
		if i != lastIndex {
			strBuilder.WriteString(sep)
		}
	}

	return strBuilder.String()
}

func JoinStr(sep string, elems ...string) string {
	return strings.Join(elems, sep)
}

//切片以，号分隔组合成字符串
func Join(elems []string) string {
	return strings.Join(elems, ",")
}

//首字母转小写
func ToFirstLetterLower(str string) string {
	chars := []rune(str)
	firstChar := chars[0]
	if unicode.IsUpper(firstChar) {
		chars[0] = unicode.ToLower(firstChar)
		return string(chars)
	} else {
		return str
	}
}

//反转字符串
func Reverse(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

// 以逗号拆分字符串
func Split(s string) []string {
	return strings.Split(s, ",")
}

// 去掉首尾空格
func Trim(s string) string {
	return strings.Trim(s, " ")
}

// interface 转 string
func ToString(value interface{}) string {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case time.Time:
		t, _ := value.(time.Time)
		key = t.String()
		// 2022-11-23 11:29:07 +0800 CST  这类格式把尾巴去掉
		key = strings.Replace(key, " +0800 CST", "", 1)
		key = strings.Replace(key, " +0000 UTC", "", 1)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// 判断是否包含有可能有注入风险的特殊符号,如果没有返回true
func SafeChars(s string) bool {
	return !strings.ContainsAny(s, UNSAFE_SQL_CHAR)
}

// 获取分隔符第一次出现之前的子字符串
func SubBefore(s string, separator string) string {
	if s == "" || separator == "" {
		return ""
	}

	i := strings.Index(s, separator)
	if i == -1 {
		return s
	}
	return s[:i]
}

// 驼峰单词转下划线单词
func CamecaseToUnderline(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
