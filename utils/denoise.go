package utils

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/dlclark/regexp2"
	"regexp"
	_const "slice-engine/const"
	"slice-engine/entity"
	"strconv"
	"strings"
	"unicode"
)

var (
	chatNoiseFullSet    = _const.CHAT_NOISE_SET
	chatNoiseTailSet    = mapset.NewSet(",", "。", "?", "？", ";", "；", "(")
	fullMatchNoiseSet   = mapset.NewSet("其它", "其他")
	chatNoisePatternSet = mapset.NewSet(regexp.MustCompile("[0-9]{17}[0-9X]"))
	digitPattern        = regexp.MustCompile("[0-9]+")

	leftBracketCompile  = "[" + strings.Join(Transform[string](_const.LEFT_BRACKET_SET.ToSlice()), "\\") + "]"
	rightBracketCompile = "[" + strings.Join(Transform[string](_const.RIGHT_BRACKET_SET.ToSlice()), "\\") + "]"
	bracketPattern      = regexp2.MustCompile(leftBracketCompile+"\\s*"+rightBracketCompile, regexp2.IgnoreCase)
)

// DeNoise removes noise from the group in the EngineState.
func DeNoise(state *entity.EngineState) {
	group := state.Group
	if len(*group) == 0 {
		return
	}
	var result []string
	for _, val := range *group {
		if strings.TrimSpace(val) == "" {
			continue
		}
		// Remove special values
		val = specValRemove(val)
		val = RemoveTopSeq(val, false)
		// Remove trailing noise characters
		for _, item := range Transform[string](chatNoiseTailSet.ToSlice()) {
			if strings.HasSuffix(val, item) {
				val = val[:len(val)-len(item)]
			}
		}
		// Check if the whole string matches any noise character
		if chatNoiseFullSet.Contains(val) {
			continue
		}
		// Check if the string is fully matched noise
		if fullMatchNoiseSet.Contains(val) {
			continue
		}
		result = append(result, val)
	}

	state.Group = &result
}

func RemoveTopSeq(text string, isStart bool) string {
	if strings.TrimSpace(text) == "" {
		return text
	}

	chars := []rune(text)
	idx := 0
	flag := false
	var ns string
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		s := string(c)
		if s == " " {
			continue
		}
		if chatNoiseFullSet.Contains(s) {
			continue
		}
		idx = i
		for unicode.IsDigit(c) {
			flag = true
			idx++
			if idx >= len(chars) {
				break
			}
			c = chars[idx]
			s = string(c)
		}
		if !flag {
			return text
		}
		ns = string(chars[i:idx])
		for c == ' ' || chatNoiseFullSet.Contains(s) {
			idx++
			if idx >= len(chars) {
				break
			}
			c = chars[idx]
			s = string(c)
		}
		if idx == 0 {
			return text
		}
		// 数字只出现了一次

		if CountMatches(text, digitPattern) > 1 {
			return text
		}
		if strings.HasPrefix(ns, "0") || len(ns) > 2 {
			return text
		}

		if n, err := strconv.Atoi(ns); err == nil {
			if n > _const.MAX_SEQ {
				return text
			}
		} else {
			return text
		}
		// 序号+符号?+[__] 可选的值包括 等级词
		sub := string(chars[idx:])
		gradeSuffix := IsGradeSuffix(sub)
		if gradeSuffix == nil || *gradeSuffix == true {
			return sub
		}
		if isStart {
			return sub
		}
		return text
	}
	return text[idx:]
}

func RemoveEscapeSymbol(text string) string {

	return text
}

func specValRemove(val string) string {
	// Replace substrings in the val using the NOISE_DIC_CACHE
	/*
		for _, noise := range NOISE_DIC_CACHE {
			val = strings.ReplaceAll(val, noise, "")
		}
	*/
	// Replace substrings matching the CHAT_NOISE_PATTERN_SET using regular expressions
	for _, item := range Transform[*regexp.Regexp](chatNoisePatternSet.ToSlice()) {
		val = item.ReplaceAllString(val, "")
	}
	// Remove empty brackets
	val = strings.TrimSpace(val)
	if len(val) == 1 {
		return val
	}
	if rs, err := bracketPattern.Replace(val, "", -1, -1); err == nil {
		val = rs
	}
	return val
}
