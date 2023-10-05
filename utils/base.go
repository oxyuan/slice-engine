package utils

import (
	"github.com/dlclark/regexp2"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func IsDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func StrIsDigit(s string) bool {
	for _, c := range s {
		if !IsDigit(c) {
			return false
		}
	}
	return true
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// CountMatches 使用正则表达式来计算一个字符串中匹配项的数量
func CountMatches(inputStr string, regex *regexp.Regexp) int {
	// 使用 FindAllString 方法获取所有匹配项
	matches := regex.FindAllString(inputStr, -1)
	// 返回匹配项的数量
	return len(matches)
}

func CountMatchesPro(inputStr string, regex *regexp2.Regexp) int {
	return -1
}

func AreSlicesEqual(slice1, slice2 any) bool {
	// 检查两个切片的类型是否相同
	if reflect.TypeOf(slice1) != reflect.TypeOf(slice2) {
		return false
	}

	// 使用反射比较两个切片的内容
	return reflect.DeepEqual(slice1, slice2)
}

// DBC2SBC 全角转半角
func DBC2SBC(s string) string {
	var strLst []string
	for _, i := range s {
		insideCode := i
		if insideCode == 12288 {
			insideCode = 32
		} else {
			insideCode -= 65248
		}
		if insideCode < 32 || insideCode > 126 {
			strLst = append(strLst, string(i))
		} else {
			strLst = append(strLst, string(insideCode))
		}
	}
	return strings.Join(strLst, "")
}

func IsAlphabeticOrDigit(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c)
}

func Transform[T any](slice []any) []T {
	if len(slice) == 0 {
		return nil
	}
	var result []T
	for _, item := range slice {
		result = append(result, item.(T))
	}
	return result
}
