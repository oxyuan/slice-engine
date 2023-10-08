package splitter

import (
	mapset "github.com/deckarep/golang-set"
	"slice-engine/entity"
	"slice-engine/utils"
	"strings"
)

var (
	HIGH_PRIORITY_SPLITTER_SET = mapset.NewSet(";", "；", "。", "?", "？", "＼")
	LOW_PRIORITY_SPLITTER_SET  = mapset.NewSet(",", "，", "、", "|")
	ALL_SPLITTER_SET           = mapset.NewSet("-")
)

func init() {
	ALL_SPLITTER_SET = ALL_SPLITTER_SET.Union(HIGH_PRIORITY_SPLITTER_SET)
	ALL_SPLITTER_SET = ALL_SPLITTER_SET.Union(LOW_PRIORITY_SPLITTER_SET)
}

func SplitOrMergeSemicolon(input string, isFuzzy bool) *[]*entity.SemicolonTuple {
	result := make([]*entity.SemicolonTuple, 0)
	arr := []rune(input)
	l := len(arr)
	position := 0
	for i := 0; i < l; i++ {
		c := arr[i]
		if c == ' ' {
			continue
		}
		if i == l-1 {
			result = append(result, &entity.SemicolonTuple{
				Text: string(arr[position:]),
				Sem:  "",
			})
			break
		}
		sem := strings.TrimSpace(string(c))

		if !isSemicolonSplitter(sem, isFuzzy) {
			continue
		}
		j := i
		for ; j+1 < l && isSemicolonSplitter(string(arr[j]), isFuzzy); j++ {
			sem = string(arr[j+1])
		}
		i = j

		_ = string(arr[position:i])
		_ = string(arr[utils.Min(i+1, l):])

		// 切分合并判断

	}
	return &result
}

func isSemicolonSplitter(sem string, isFuzzy bool) bool {
	b := LOW_PRIORITY_SPLITTER_SET.Contains(sem)
	if !b {
		b = b || HIGH_PRIORITY_SPLITTER_SET.Contains(sem)
	}
	return b
}

func isFirstMergeMatch() {

}
