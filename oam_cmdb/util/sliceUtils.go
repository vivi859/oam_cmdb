package util

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/slices"
)

type IntList []int
type StrList []string

// 提供一个元素类型转换函数mapper，将传入的切片各元素转为其他类型，并组成新的切片返回
func SliceConvert[T any, R any](slice []T, mapper func(T) R) []R {
	if len(slice) == 0 {
		return nil
	}
	tmp := make([]R, 0)
	for _, ele := range slice {
		tmp = append(tmp, mapper(ele))
	}
	return tmp
}

// 从切片中筛选满足给定条件的元素组成新的切片，并返回
func SliceFilter[T any](slice []T, filter func(T) bool) []T {
	if len(slice) == 0 {
		return nil
	}
	tmp := make([]T, 0)
	for _, ele := range slice {
		if filter(ele) {
			tmp = append(tmp, ele)
		}
	}
	return tmp
}

// 从切片中筛选满足条件的第一个元素
func SliceFilterOne[T any](slice []T, filter func(T) bool) (bool, T) {
	var result T
	if len(slice) > 0 {
		for _, ele := range slice {
			if filter(ele) {
				return true, ele
			}
		}
	}
	return false, result
}

func (list IntList) ToJSONString() string {
	str, _ := jsoniter.MarshalToString(list)
	return str
}

//判断元素是否在集合中
func (list IntList) IsContain(ele int) bool {
	for _, vv := range list {
		if vv == ele {
			return true
		}
	}
	return false
}

//将集合元素合并为字符串,sep为分隔符
func (list IntList) Join(sep string) string {
	return JoinInteger(sep, list...)
}

// 字符串切片转为整形切片,必须保证切片元素都是可转换的,否则引发panic
func ToIntSlice(elems []string) []int {
	return SliceConvert[string, int](elems, func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return i
	})

}

// 将切片转为元素是interface{}类型的切片
func ToInterfaceSlice[T interface{}](elems []T) []interface{} {
	arr := make([]interface{}, len(elems))
	for i := range elems {
		arr[i] = elems[i]
	}
	return arr
}

func ToStrSlice(elems []interface{}) []string {
	arr := make([]string, len(elems))
	for i := range elems {
		arr[i] = elems[i].(string)
	}
	return arr
}

// 删除切片中指定元素
func RemoveSlice[T comparable](slice []T, elem T) []T {
	if len(slice) == 0 {
		return slice
	}
	tmpSlice := slice[:0]
	for _, v := range slice {
		if v != elem {
			tmpSlice = append(tmpSlice, v)
		}
	}

	return tmpSlice
}

func MapToIntSlice[T any](slice []T, mapper func(src T) int) []int {
	return SliceConvert[T, int](slice, mapper)
}

//合并两个切片，重复元素只保留一个
func UnionSlice[T comparable](s1, s2 []T) []T {
	var newslice []T
	if len(s1) > 0 {
		if len(s2) > 0 {
			newslice = append(newslice, s1...)
			for _, e := range s2 {
				if !slices.Contains[T](newslice, e) {
					newslice = append(newslice, e)
				}
			}
		} else {
			newslice = s1
		}
	} else {
		if len(s2) > 0 {
			newslice = s2
		}
	}

	return newslice
}

//切片转map
func SliceToMap[S any, K comparable, V any](slice []S, keyMapper func(S) K, valueMapper func(S) V) map[K]V {
	size := len(slice)
	if size > 0 {
		newmap := make(map[K]V, size)
		for _, ele := range slice {
			k := keyMapper(ele)
			newmap[k] = valueMapper(ele)
		}
		return newmap
	}
	return nil
}
