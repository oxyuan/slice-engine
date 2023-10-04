package engine

import (
	"github.com/rs/zerolog/log"
	"math"
	"regexp"
	_const "slice-engine/const"
	"slice-engine/entity"
	"slice-engine/enums"
	"slice-engine/splitter"
	"slice-engine/utils"
	"strconv"
	"strings"
)

var seqStartPattern = regexp.MustCompile("^[\\s<\\[\\(]*(\\d+)")

type OrdinalEngine struct {
	RootEngine
}

func init() {
	FACTORY = append(FACTORY, &OrdinalEngine{})
}

func (e *OrdinalEngine) GetCategory() enums.EngineType {
	return enums.EngineType_ORDINAL
}

func (e *OrdinalEngine) Recognize(state *entity.EngineState) bool {
	text := state.Text
	// if the text is empty, return false.
	if text == "" || len(text) == 0 {
		return false
	}
	// 文本序号切割
	tuple := splitter.Split(text)
	state.OrdinalTuple = &tuple

	seqList := make([]float64, 0)
	for _, t := range tuple {
		if t.Seq != _const.DEF_SEQ {
			seqList = append(seqList, t.Seq)
		}
	}
	size := len(seqList)
	if size == 0 {
		return false
	}
	if size == 1 {
		seq := seqList[0]
		return seq > 0 && seq <= 1
	}

	// seqList 的值如果有连续三个递增，返回true
	count := 0
	for i := 1; i < size-1; i++ {
		if seqList[i-1]+1 == seqList[i] {
			count++
		} else {
			count = 0
		}
		if count >= 2 {
			return true
		}
	}

	// 判断 seqList 的排序值是否全是逆序
	reverse := true
	for i := len(seqList) - 1; i >= 1; i-- {
		seq1 := seqList[i-1]
		seq2 := seqList[i]
		if seq1 < seq2 {
			reverse = false
			break
		}
	}
	if reverse {
		return false
	}

	// 判断匹配的数字是否是序号
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, seq := range seqList {
		if seq == _const.DEF_SEQ {
			continue
		}
		min = math.Min(min, seq)
		max = math.Max(max, seq)
	}

	return (max-min)/float64(size) < _const.DEF_SEQ_INCREMENT_RATE
}

func (e *OrdinalEngine) Slice(state *entity.EngineState) {
	group := make([]string, 0)
	for _, tuple := range *state.OrdinalTuple {
		group = append(group, tuple.Text)
	}
	state.Group = &group
}

func (e *OrdinalEngine) Merge(state *entity.EngineState) {
	tupleList := state.OrdinalTuple
	mergeList := make([]*[]string, 0)

	visited := make([][]bool, 2)
	for i := range visited {
		visited[i] = make([]bool, len(*tupleList))
	}

	backtrack(&mergeList, tupleList, &[]string{}, visited, 0)

	state.MergeList = &mergeList
}

func backtrack(mergeList *[]*[]string, tuples *[]*entity.OrdinalTuple, branch *[]string, visited [][]bool, index int) {
	if len(*mergeList) > 2<<10 {
		msg := "the merge list is too large. group: ["
		for _, tuple := range *tuples {
			msg += tuple.Text + " "
		}
		msg += "]"

		log.Error().Msg(msg)
		return
	}
	if len(*branch) > len(*tuples) {
		return
	}
	if index == len(*tuples) && isFullTrue(visited) {
		newBranch := make([]string, len(*branch))
		copy(newBranch, *branch)
		*mergeList = append(*mergeList, &newBranch)
		return
	}

	for i := index; i < len(*tuples); i++ {
		tuple := (*tuples)[i]
		group := tuple.Text
		seq := tuple.Seq
		if isTerminate(visited, i) {
			return
		}
		if len(*branch) == 0 || isFree(tuples, *branch, index, seq) {
			*branch = append(*branch, group)
			visited[0][i] = true
			backtrack(mergeList, tuples, branch, visited, index+1)
			*branch = (*branch)[:len(*branch)-1]
			visited[0][i] = false
		}
		if i == 0 || isFullFalse(visited) || isTerminate(visited, i) {
			return
		}

		// 剪枝
		if prune(tuples, *branch, i) {
			return
		}

		visited[1][i] = true
		appendToBranch(tuples, branch, i)
		backtrack(mergeList, tuples, branch, visited, index+1)
		substringFromBranch(tuples, branch, i)
		visited[1][i] = false
	}
}

func (e *OrdinalEngine) Optimize(state *entity.EngineState) {
	tupleList := state.OrdinalTuple
	mergeList := state.MergeList
	if len(*tupleList) == 0 {
		return
	}
	priorityMap := make(map[int]*entity.OrdinalScore)
	suffix := seqSub(tupleList)
	idx := 0
	for i, group := range *mergeList {
		score := float64(len(*group))

		// 序号后缀是否统一
		incBool := false
		for j, text := range *group {
			if len(text) <= 1 {
				score--
				continue
			}
			if suffix != getSeqSuffix(text) && !incBool {
				continue
			}
			incBool = true
			if j == 0 {
				suffix = getSeqSuffix(text)
				continue
			}
			if suffix != getSeqSuffix(text) {
				score--
			}
		}

		incSeq := float64(0)
		prevSeq := float64(0)
		incBool = false
		for j, text := range *group {
			seq := getSeq(text)
			if seq == _const.DEF_SEQ {
				continue
			}
			if !incBool {
				incSeq = seq
				prevSeq = seq
			}
			if j == 0 {
				continue
			}
			incBool = true

			span := math.Min(seq-incSeq, seq-prevSeq)
			span = math.Min(span, 2)

			if span == 1 {
				incSeq = seq
			} else {
				score -= span
			}
			prevSeq = seq
		}
		priorityMap[i] = &entity.OrdinalScore{
			Score:     score,
			IsPerfect: score == float64(len(*group)),
		}
	}
	score := float64(math.MinInt32)
	for key, value := range priorityMap {
		isBetter := value.Score > score
		isFull := value.IsPerfect && value.Score == score
		if isBetter || isFull {
			idx = key
			score = value.Score
		}
	}
	state.Group = (*mergeList)[idx]
}

func (e *OrdinalEngine) Standardize(state *entity.EngineState) {
	utils.Standardize(state)
}

func (e *OrdinalEngine) DeNoise(state *entity.EngineState) {
	removeSeq(state)
	utils.DeNoise(state)
}

func prune(tuples *[]*entity.OrdinalTuple, branch []string, i int) bool {
	size := len(branch)
	prevVal := branch[size-1]
	return pruneWithPrevVal(tuples, prevVal, i)
}

func pruneWithPrevVal(tuples *[]*entity.OrdinalTuple, prevVal string, i int) bool {
	tuple := (*tuples)[i]

	prevValSeqSuffix := getSeqSuffix(prevVal)
	prevValSeq := getSeq(prevVal)
	currVal := tuple.Text
	currValSeqSuffix := getSeqSuffix(currVal)
	currValSeq := getSeq(currVal)

	isCurValGreaterThanPrevVal := currValSeq-prevValSeq == 1

	// 1. 如果自身与待拼接的文本的 [序号=pre+1] && [序号后缀一致]，不能拼接
	isEqualSuffix := prevValSeqSuffix == currValSeqSuffix
	if isCurValGreaterThanPrevVal && isEqualSuffix {
		return true
	}

	// 2. 如果自身与待拼接的文本的 [序号=pre+1] && [序号=next-1]，不能拼接
	isLargeEnough := currValSeq >= _const.ENOUGH_LARGE_SEQ && len(*tuples) > _const.ENOUGH_LARGE_SEQ
	if isCurValGreaterThanPrevVal && isLargeEnough {
		return true
	}

	// 3. 如果自身与待拼接的文本的 [序号=pre+1] && [序号已经足够大]，不能拼接
	isNextValGreaterThanCurVal := false
	if i+1 != len(*tuples) {
		nextVal := (*tuples)[i+1].Text
		nextValSeq := getSeq(nextVal)
		isNextValGreaterThanCurVal = nextValSeq-currValSeq == 1
	}
	if isCurValGreaterThanPrevVal && isNextValGreaterThanCurVal {
		return true
	}

	// 4. 如果自身与待拼接的文本的 [序号=pre.pre+1] && [序号=next-1] && [序号已经足够大]，不能拼接
	isPprevValGreaterThanPrevVal := false
	if isLargeEnough && isNextValGreaterThanCurVal && i >= 2 {
		pprevValSeq := getSeq((*tuples)[i-2].Text)
		isPprevValGreaterThanPrevVal = currValSeq-pprevValSeq == 1
		if isPprevValGreaterThanPrevVal {
			return true
		}
	}

	// 5. 待续

	return false
}

func seqSub(tuples *[]*entity.OrdinalTuple) (prefix string) {
	prefix = ""
	if len(*tuples) == 0 {
		return
	}
	prefixMap := make(map[string]int)
	maxTimes := 0
	root := (*tuples)[0].Suffix
	for _, tuple := range *tuples {
		s := (*tuple).Suffix
		times := 0
		if times, ok := prefixMap[s]; ok {
			prefixMap[s] = times + 1
		} else {
			prefixMap[s] = 1
		}
		if times > maxTimes {
			prefix = s
			maxTimes = times
		} else if times == maxTimes && s == root {
			prefix = s
		}
	}
	return
}

func getSeqSuffix(text string) (suffix string) {
	suffix = ""
	if len(text) == 0 || text == "" {
		return
	}
	// 遍历文本，获取序号后的第一个值
	point := false
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if c == ' ' {
			continue
		}
		if c >= 48 && c <= 57 {
			point = true
			continue
		}
		if !point && !_const.CHAT_NOISE_SET.Contains(string(c)) {
			return
		}
		s := string(c)
		if point {
			if _const.CHAT_NOISE_SET.Contains(s) {
				return s
			}
			return
		}
	}
	return
}

func getSeq(text string) float64 {
	if len(text) == 0 || text == "" {
		return _const.DEF_SEQ
	}
	seqList := seqStartPattern.FindAllString(text, 1)
	if len(seqList) == 0 || len(seqList[0]) > 0 {
		return _const.DEF_SEQ
	}
	s := seqList[0]
	if seq, err := strconv.Atoi(s); err != nil || seq < 0 || seq > 2<<4 {
		return _const.DEF_SEQ
	} else {
		return float64(seq)
	}
}

func isFullTrue(visited [][]bool) bool {
	for i := 0; i < len(visited[0]); i++ {
		count := 0
		for _, booleans := range visited {
			if booleans[i] {
				count++
			}
		}
		if count != 1 {
			return false
		}
	}
	return true
}

func isTerminate(visited [][]bool, i int) bool {
	// 当前分支已经被选择（拼接 OR 添加）
	if visited[0][i] || visited[1][i] {
		return true
	}

	// 不能跨节点选择，当前分支异常
	if i != 0 && !visited[0][i-1] && !visited[1][i-1] {
		return true
	}

	// 待续
	return false
}

func isFree(tuples *[]*entity.OrdinalTuple, branch []string, index int, seq float64) bool {
	root := (*tuples)[0].Seq

	// 是否不存在序号
	isWithoutSeq := seq == -1
	if isWithoutSeq {
		// return false;
	}

	// 序号是否乱序
	isOutOfOrder := seq == -1 || seq == root
	if isOutOfOrder {
		return false
	}

	// 序号是否和前一个重复
	preSeq := -1.0
	if index > 0 {
		preSeq = getSeq(branch[len(branch)-1])
	}
	isRepeat := seq == preSeq
	if isRepeat {
		return false
	}

	//	// 当前序号是否是一个单位
	//	pre := ""
	//	if index > 0 {
	//		pre = groups[index-1]
	//	}
	//	cur := groups[index]
	//	next := ""
	//	if index < len(groups)-1 {
	//		next = groups[index+1]
	//	}
	//	isUnit := isUnit(pre, cur, next)
	//	if isUnit {
	//		return false
	//	}

	// todo
	return true
}

func isFullFalse(visited [][]bool) bool {
	for _, booleans := range visited {
		for _, boolValue := range booleans {
			if boolValue {
				return false
			}
		}
	}
	return true
}

func appendToBranch(tuples *[]*entity.OrdinalTuple, branch *[]string, index int) {
	if len(*branch) == 0 {
		return
	}
	last := (*branch)[len(*branch)-1]
	(*branch)[len(*branch)-1] = last + (*tuples)[index].Text
}

func substringFromBranch(tuples *[]*entity.OrdinalTuple, branch *[]string, index int) {
	if len(*branch) == 0 {
		return
	}
	group := (*tuples)[index].Text
	last := (*branch)[len(*branch)-1]
	if len(last) < len(group) {
		return
	}
	(*branch)[len(*branch)-1] = last[:len(last)-len(group)]
}

func removeSeq(state *entity.EngineState) {
	group := make([]string, 0)
	for _, text := range *state.Group {
		suffix := getSeqSuffix(text)
		chars := []rune(text)
		index := 0
		for _, c := range chars {
			if c == ' ' || c >= 48 && c <= 57 {
				index++
				continue
			}
			if len(strings.TrimSpace(suffix)) == 0 {
				break
			}
			if c == rune(suffix[0]) {
				index++
				break
			}
		}
		group = append(group, string(chars[index:]))
	}
	state.Group = &group
}
