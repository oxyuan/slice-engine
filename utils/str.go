package utils

import (
	_const "slice-engine/const"
	"strings"
	"unicode"
)

var FALSE = false
var TRUE = true

func IsSerialNumberOfEndAndStart(prefix, suffix string, freq int) bool {
	if prefix == "" || suffix == "" {
		return false
	}

	prefixChars := []rune(prefix)
	suffixChars := []rune(suffix)
	times := freq

	for i := len(prefixChars) - 1; i >= 0; i-- {
		c := prefixChars[i]
		if isContinueChar(c) {
			continue
		}
		if unicode.IsDigit(c) {
			times--
			continue
		}
		if times <= 0 {
			return true
		}
		break
	}

	times = freq

	for i := 0; i < len(suffixChars); i++ {
		c := suffixChars[i]
		if isContinueChar(c) {
			continue
		}
		if unicode.IsDigit(c) {
			times--
			continue
		}
		if times <= 0 {
			return true
		}
		break
	}

	return false
}

func isContinueChar(c rune) bool {
	if c == ' ' {
		return true
	}
	s := string(c)
	return _const.CHAT_NOISE_SET.Contains(s)
}

func IsContinueChar(c rune) bool {
	if c == ' ' {
		return true
	}
	s := string(c)
	return _const.CHAT_NOISE_SET.Contains(s)
}

func IsGradeSuffix(val string) *bool {
	chars := []rune(val)
	return isGradeSuffixFromChars(chars)
}

func isGradeSuffixFromChars(chars []rune) *bool {
	isGradeSuffix := false
	for i := 0; i < len(chars); i++ {
		c := chars[i]
		if c == ' ' {
			continue
		}
		s := string(c)
		if IsAlphabeticOrDigit(c) || _const.STAGE_SET.Contains(s) {
			isGradeSuffix = true
			continue
		}
		if _const.XINA_SET.Contains(s) || _const.CONNECTOR_SET.Contains(s) {
			continue
		}
		if isGradeSuffix {
			if _const.GRADE_SUFFIX_SET.Contains(s) {
				if i == len(chars)-1 {
					return &TRUE
				}
				next := chars[i+1]
				isGradeSuffix = next == ' ' || _const.CHAT_NOISE_SET.Contains(string(next))
				if isGradeSuffix {
					return &TRUE
				}
				return nil
			}
			sub := string(chars[i:])
			for _, word := range Transform[string](_const.DESCRIPTIVE_WORD_SET.ToSlice()) {
				if strings.HasPrefix(sub, word) {
					return &TRUE
				}
			}
		}
		break
	}
	return &FALSE
}
