package utils

import (
	"github.com/dlclark/regexp2"
	_const "slice-engine/const"
	"slice-engine/entity"
	"strings"
)

var (
	sequenceMap = map[string]*regexp2.Regexp{
		"1" + seqSuffix: regexp2.MustCompile(`^[1①⑴⒈ⅰⅠ][\\.，。；、]`, regexp2.IgnoreCase),
		"2" + seqSuffix: regexp2.MustCompile(`^[2②⒉⑵ⅱⅡ][\\.，。；、]`, regexp2.IgnoreCase),
		"3" + seqSuffix: regexp2.MustCompile(`^[3⒊⑶③ⅲⅢ][\\.，。；、]`, regexp2.IgnoreCase),
		"4" + seqSuffix: regexp2.MustCompile(`^[4④⒋⑷㈣ⅳⅣ][\\.，。、；]`, regexp2.IgnoreCase),
		"5" + seqSuffix: regexp2.MustCompile(`^[5⑤㈤⒌⑸Ⅴⅴ][\\.，。、；]`, regexp2.IgnoreCase),
		"6" + seqSuffix: regexp2.MustCompile(`^[6⑥㈥⑹ⅵ⒍Ⅵ][\\.，。、；]`, regexp2.IgnoreCase),
		"7" + seqSuffix: regexp2.MustCompile(`^[7㈦⒎ⅶⅦ㈦⑦][\\.，。、；]`, regexp2.IgnoreCase),
		seqSuffix:       regexp2.MustCompile(`(?<=^[a-zA-Z]{1,2})[、,，\\.]`, regexp2.IgnoreCase),
	}
	replaceMap               = map[string]string{"【": "[", "】": "]", "（": "(", "）": ")"}
	seqSuffix                = "、"
	gradeSuffixPatternString = strings.Join(Transform[string](_const.CONNECTOR_SET.ToSlice()), "")
	patternReplaceTupleList  = []entity.PatternReplaceTuple{
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])[1①⑴⒈ⅰI]\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅰ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:2|②|⒉|⑵|ⅱ|11|ii|II)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅱ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:3|③|⒊|⑶|ⅲ|111|iii|III)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅲ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:4|④|⒋|⑷|ⅳ|iv|IV)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅳ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:5|⑤|⒌|⑸|ⅴ|v|V)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅴ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:6|⑥|⒍|⑹|ⅵ|vi|VI)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅵ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:7|⑦|⒎|⑺|ⅶ|vii|VII)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅶ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:8|⑧|⒏|⑻|ⅷ|viii|VIII)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅷ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:9|⑨|⒐|⑼|ⅸ|ix|IX)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅸ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:10|⑩|⒑|⑽|ⅹ|x|X)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅹ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:11|⑪|⒒|⑾|ⅺ|xi|XI)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅺ",
		},
		{
			Matter:  regexp2.MustCompile(`(?<![a-zA-Z])(?:12|⑫|⒓|⑿|ⅻ|xii|XII)\s*(?=[`+gradeSuffixPatternString+`])`, regexp2.IgnoreCase),
			Replace: "Ⅻ",
		},
	}
)

func PreStandard(text string) string {
	if text == "" || len(text) == 0 {
		return text
	}
	for _, tuple := range patternReplaceTupleList {
		text, _ = tuple.Matter.Replace(text, tuple.Replace, -1, -1)
	}
	for k, v := range replaceMap {
		text = strings.ReplaceAll(text, k, v)
	}
	return text
}

func Standardize(state *entity.EngineState) {
	group := state.Group
	if group == nil || len(*group) == 0 {
		return
	}
	result := make([]string, 0)
	for _, item := range *group {
		item = strings.TrimSpace(item)
		// 序号前缀去除
		item = removePrefix(item)
		// 序号标化
		item = sequence(item)
		// 序号后的罗马数字标化
		item = romanNumConv(item)
		result = append(result, item)
	}
	state.Group = &result
}

// 序号前缀去除
func removePrefix(val string) string {
	var sb strings.Builder
	chars := []rune(val)
	for i := 0; i < Min(len(chars), 3); i++ {
		c := chars[i]
		if c == ' ' || _const.CHAT_NOISE_SET.Contains(string(c)) {
			continue
		}
		sb.WriteString(string(chars[i:]))
		break
	}
	return sb.String()
}

func sequence(val string) string {
	for standardVal, pattern := range sequenceMap {
		if rs, err := pattern.Replace(val, standardVal, 0, 1); err == nil && rs != val {
			return rs
		}
	}
	return val
}

func romanNumConv(val string) string {
	split := strings.Split(val, seqSuffix)
	length := len(split)
	if length < 2 {
		return val
	}
	s1 := split[0]
	var c strings.Builder
	for i := 1; i < len(split); i++ {
		c.WriteString(split[i])
		if i != len(split)-1 {
			c.WriteString(seqSuffix)
		}
	}
	s2 := c.String()
	return s1 + seqSuffix + s2
}
