package utils

import (
	"regexp"
	_const "slice-engine/const"
	"strings"
)

func Slice(text string, regexp *regexp.Regexp) *[]string {
	groups := make([]string, 0)
	splits := make([]string, 0)
	for _, s := range regexp.Split(text, -1) {
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		splits = append(splits, s)
	}
	seqList := regexp.FindAllString(text, -1)
	length := len(splits)
	size := len(seqList)
	if length == 0 || size == 0 || length < size {
		return &[]string{text}
	}
	split0 := splits[0]
	if length == size && split0 == "" {
		for i := 1; i < len(seqList); i++ {
			groups = append(groups, seqList[i-1]+splits[i])
		}
		groups = append(groups, seqList[len(seqList)-1])
	} else {
		groups = append(groups, split0)
		for i := 0; i < len(seqList); i++ {
			val := ""
			if i+1 != len(splits) {
				val = splits[i+1]
				groups = append(groups, seqList[i]+val)
			}
		}
	}
	if groups[len(groups)-1] == seqList[len(seqList)-1] && len(groups) > 1 {
		s := groups[len(groups)-2]
		// 移除 groups 中的最后一个元素
		groups = groups[:len(groups)-1]
		s = groups[len(groups)-1] + s
		groups = groups[:len(groups)-1]
		groups = append(groups, s)
	}
	if len(groups) <= 1 {
		return &groups
	}
	distinct := Distinct(&groups)
	// 移除 distinct 中所有的空字符串
	for i := 0; i < len(*distinct); i++ {
		if (*distinct)[i] == "" {
			*distinct = append((*distinct)[:i], (*distinct)[i+1:]...)
			i--
		}
	}
	return distinct
}

func Distinct(origin *[]string) *[]string {
	groups := make([]string, len(*origin))
	for i, group := range *origin {
		groups[i] = strings.TrimSpace(group)
	}

	if len(groups) <= 1 {
		return origin
	}

	distinctMap := make(map[string]struct{})
	distinctGroups := make([]string, 0)

	for _, group := range groups {
		if _, ok := distinctMap[group]; !ok {
			distinctMap[group] = struct{}{}
			distinctGroups = append(distinctGroups, group)
		}
	}

	if len(distinctGroups) == len(groups) {
		return origin
	}

	// 全部重复
	repeatIndex := make([]int, 0)
	for i := 1; i < len(groups); i++ {
		group := groups[i]
		idx := findIndex(group, groups)

		if idx == i && idx < len(distinctGroups) {
			continue
		}
		if isRepeat(i, repeatIndex, groups) {
			repeatIndex = append(repeatIndex, i)
		}
	}

	sum := 0
	for i := 1; i < len(repeatIndex); i++ {
		sum += repeatIndex[i] - repeatIndex[i-1]
	}

	if len(repeatIndex) > 0 && float64(sum/len(repeatIndex)) < _const.DESTINCT_INCREASE_RATE {
		end := repeatIndex[0]
		newSlice := (*origin)[:end]
		return &newSlice
	}
	return origin
}

func findIndex(group string, groups []string) int {
	for i, g := range groups {
		if g == group {
			return i
		}
	}
	return -1
}

func isRepeat(i int, repeatIndex []int, groups []string) bool {
	lastIdx := len(repeatIndex) - 1
	group := groups[i]

	// 完全重复
	idx := findIndex(group, groups)
	if idx != -1 && idx != i {
		return true
	}

	// 后缀重复
	if strings.HasSuffix(groups[lastIdx], group) {
		return true
	}

	// 前缀重复
	if strings.HasPrefix(groups[lastIdx], group) {
		return true
	}

	// 待补充
	return false
}
