package splitter

import (
	mapset "github.com/deckarep/golang-set"
	"regexp"
	_const "slice-engine/const"
	"slice-engine/entity"
	"slice-engine/utils"
	"strconv"
	"strings"
	"unicode"
)

var (
	// illegalPrefixSet 非法前缀
	illegalPrefixSet = mapset.NewSet("-", "+", "_", "第", "尾")
	// illegalPatternPrefixSet 非法前缀（多值组合-正则）
	illegalPatternPrefixSet = mapset.NewSet(
		regexp.MustCompile("孕(\\d+产)?$"),
		regexp.MustCompile("\\d[\\.、/]$"),
	)
	// illegalSuffixSet 非法后缀（单位-字符串）
	illegalSuffixSet = _const.UNIT_SUFFIX_SET
	// illegalPatternSuffixSet 非法后缀（多值组合-正则）
	illegalPatternSuffixSet = mapset.NewSet(
		regexp.MustCompile("^(loa|lot|lop|roa|rot|rop)"),
		regexp.MustCompile("^\\s*[\\.\\*\\^\\-\\+/][0-9]+"),
	)
	// illegalTupleSet 当某个前缀和后缀同时出现时，禁止
	illegalTupleSet = mapset.NewSet()
	// illegalPatternTupleSet 当某个前缀和后缀同时出现时，禁止（多值组合-正则）
	illegalPatternTupleSet = mapset.NewSet(
		entity.NewPreSufTupleWithRegexp(
			regexp.MustCompile("[\\(\\[（【][\\w\\.；]*$"),
			regexp.MustCompile("^[\\w\\.;；]*[\\)\\]）】]")),
	)
	// allowedTupleSet 合法前缀
	allowedTupleSet = mapset.NewSet(
		entity.NewPreSufTupleWithString("(", ")"),
		entity.NewPreSufTupleWithString("（", "）"),
		entity.NewPreSufTupleWithString("[", "]"),
	)
)

func Split(input string) []*entity.OrdinalTuple {
	texts := make([]string, 0)
	seqs := make([]string, 0)

	position := 0

	runes := []rune(input)
	size := len(runes)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if !unicode.IsDigit(c) && i != size-1 {
			continue
		}
		// catch the number
		j := i
		for j+1 < size && unicode.IsDigit(runes[j+1]) {
			j++
		}

		prefix := string(runes[position:i])
		suffix := string(runes[utils.Min(j+1, size):size])

		if i == size-1 {
			texts = append(texts, string(runes[position:size]))
			break
		}

		// 合法分词位点
		if isAllowed(prefix, suffix) {
			texts = append(texts, string(runes[position:i]))
			seqs = append(seqs, string(runes[i:j+1]))
			position = i
			i = j
			continue
		}
		// 非法分词位点
		if isIllegal(prefix, suffix) {
			i = j
			continue
		}
		if position != i {
			texts = append(texts, string(runes[position:i]))
		}
		seqs = append(seqs, string(runes[i:j+1]))
		position = i
		i = j
	}
	texts = *utils.Distinct(&texts)

	seqLen := len(seqs)
	textLen := len(texts)

	j := textLen - seqLen

	ordinalTuple := make([]*entity.OrdinalTuple, 0)
	for i, text := range texts {
		seq := ""
		if seqLen != 0 && i >= j {
			seq = seqs[i-j]
		}
		assignTuple(&ordinalTuple, text, seq)
	}
	return ordinalTuple
}

func assignTuple(list *[]*entity.OrdinalTuple, text, seq string) {
	if text == "" {
		return
	}
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.TrimSpace(text)

	suffix := ""
	position := 0

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		if char == ' ' {
			continue
		}
		val := string(char)
		if _const.CHAT_NOISE_SET.Contains(val) {
			suffix = val
			position = i + 1
			break
		}
	}
	*list = append(*list, &entity.OrdinalTuple{
		Text:     text,
		Seq:      getSeq(seq),
		Suffix:   suffix,
		Position: position,
	})
}

func isAllowed(prefix, suffix string) bool {
	if utils.AnySetMatch[entity.PreSufTuple](allowedTupleSet, func(val entity.PreSufTuple) bool {
		return strings.HasSuffix(prefix, val.Pre) && strings.HasPrefix(suffix, val.Suf)
	}) {
		return true
	}

	return false
}

func isIllegal(prefix, suffix string) bool {
	if prefix == "" || suffix == "" {
		return false
	}
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	suffix = strings.ToLower(strings.TrimSpace(suffix))

	if utils.AnySetMatch[string](illegalPrefixSet, func(val string) bool {
		return strings.HasSuffix(prefix, val)
	}) {
		return true
	}

	if utils.AnySetMatch[string](illegalSuffixSet, func(val string) bool {
		return strings.HasPrefix(suffix, val)
	}) {
		return true
	}

	if utils.AnySetMatch[entity.PreSufTuple](illegalTupleSet, func(val entity.PreSufTuple) bool {
		return strings.HasSuffix(prefix, val.Pre) && strings.HasPrefix(suffix, val.Suf)
	}) {
		return true
	}

	if utils.AnySetMatch[*regexp.Regexp](illegalPatternPrefixSet, func(val *regexp.Regexp) bool {
		return val.MatchString(prefix)
	}) {
		return true
	}

	if utils.AnySetMatch[*regexp.Regexp](illegalPatternSuffixSet, func(val *regexp.Regexp) bool {
		return val.MatchString(suffix)
	}) {
		return true
	}

	if utils.AnySetMatch[entity.PreSufTuple](illegalPatternTupleSet, func(val entity.PreSufTuple) bool {
		return (val.PreRegexp != nil && val.PreRegexp.MatchString(prefix)) && (val.SufRegexp != nil && val.SufRegexp.MatchString(suffix))
	}) {
		return true
	}

	// 前后有字母时
	if hasLetter(prefix, suffix) {
		return true
	}

	// 前后的值是连续的序号组成时
	if utils.IsSerialNumberOfEndAndStart(prefix, suffix, 2) {
		return true
	}

	return false
}

func hasLetter(prefix, suffix string) bool {
	if strings.TrimSpace(prefix) == "" || strings.TrimSpace(suffix) == "" {
		return false
	}
	c := prefix[len(prefix)-1]
	p := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
	c = suffix[0]
	s := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
	return p && s
}

func getSeq(seq string) float64 {
	if strings.TrimSpace(seq) == "" {
		return _const.DEF_SEQ
	}

	if _, err := strconv.ParseFloat(seq, 64); err != nil {
		return _const.DEF_SEQ
	}

	seq = strings.TrimSpace(seq)
	v, err := strconv.ParseFloat(seq, 64)
	if err != nil {
		return _const.DEF_SEQ
	}

	if v < 0 || v > _const.MAX_SEQ {
		return _const.DEF_SEQ
	}

	return v
}
