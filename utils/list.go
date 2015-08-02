package utils

import (
	"sort"
)

//列表是否包含给定项
func ListContains(list []interface{}, key interface{}) (finded bool) {

	for _, v := range list {
		if v == key {
			finded = true
			break
		}
	}
	return
}

//字符串数组中是否包含给定项
func StringsContains(list []string, key string) (finded bool) {
	for _, v := range list {
		if v == key {
			finded = true
			break
		}
	}
	return
}

//int数组中是否包含给定项
func Int64sContains(list []int64, key int64) (finded bool, index int) {
	for i, v := range list {
		if v == key {
			finded = true
			index = i
			break
		}
	}
	return
}

//移除slice中的元素
func RemoveStringSlice(slice []string, s string) []string {
	sort.Strings(slice)
	i := sort.SearchStrings(slice, s)

	if i < len(slice) && slice[i] == s {
		return append(slice[:i], slice[i+1:]...)
	}
	return slice
}

//移除slice中的元素
func RemoveInt64Slice(slice []int64, s int64) []int64 {
	_, i := Int64sContains(slice, s)

	if i < len(slice) && slice[i] == s {
		return append(slice[:i], slice[i+1:]...)
	}
	return slice
}

// int64型slice转string型slice
func Int64s2Strings(i64 []int64) []string {
	s := make([]string, len(i64))
	for index, i := range i64 {
		s[index] = Int642str(i)
	}

	return s
}
