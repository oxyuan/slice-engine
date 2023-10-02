package boot

import (
	mapset "github.com/deckarep/golang-set"
	_const "slice-engine/const"
	"slice-engine/engine"
	"slice-engine/entity"
	"slice-engine/utils"
	"sort"
	"strings"
)

func init() {
	sortByCategory(engine.FACTORY)
}

func Run(val string) *[]string {
	// 预处理
	val = preTreat(val)

	state := entity.NewEngineState(val)
	group := new([]string)
	loop(state, group, 0)

	// 后处理
	group = afterTreat(group)

	return group
}

func loop(state *entity.EngineState, result *[]string, depth uint8) {
	text := state.Text
	origin := make([]string, 0)
	if state.Group == nil || len(*state.Group) == 0 {
		origin = append(origin, text)
	} else {
		origin = *state.Group
	}
	if text == "" || len(text) == 0 {
		return
	}
	// bootstrap
	execute(state)

	if len(*state.Group) == 0 {
		*result = append(*result, text)
		return
	}

	if state.IsBreak() {
		state.Text = (*state.Group)[0]
		// 二轮识别
		if depth == 0 {
			state.IsFuzzy = true
			loop(state, result, depth+1)
			return
		}
		*result = append(*result, *state.Group...)
		return
	}
	if utils.AreSlicesEqual(origin, *state.Group) {
		*result = append(*result, *state.Group...)
		return
	}

	for _, item := range *state.Group {
		state.Text = item
		loop(state, result, depth+1)
	}
}

func execute(state *entity.EngineState) {
	state.IsRecognized = false

	for _, e := range engine.FACTORY {
		if state.IsExcluded(&e) {
			continue
		}

		r := engine.RootEngine{
			Engine: e,
		}
		r.Run(state)

		e.Standardize(state)
		e.DeNoise(state)

		if !state.IsRecognized {
			continue
		}
		break
	}
}

func preTreat(text string) string {
	if text == "" || len(text) == 0 {
		return text
	}
	// 去除特殊符号
	text = conjointTreat(text)
	// 预标化
	text = utils.PreStandard(text)
	// 前置序号去除
	text = utils.RemoveTopSeq(text, true)
	// 转义符号去除
	text = utils.RemoveEscapeSymbol(text)

	return text
}

func afterTreat(group *[]string) (result *[]string) {
	result = new([]string)
	if group == nil || len(*group) == 0 {
		return
	}
	distinct := mapset.NewSet()
	for _, item := range *group {
		if item == "" || len(item) == 0 {
			continue
		}
		// 全半角转化
		item = utils.DBC2SBC(item)
		item = conjointTreat(item)
		// 去除空格
		item = strings.ReplaceAll(item, " ", "")
		if item != "" && len(item) != 0 {
			if distinct.Contains(item) {
				continue
			}
			*result = append(*result, item)
			distinct.Add(item)
		}
	}
	return
}

func sortByCategory(engines []entity.Engine) {
	sort.Slice(engines, func(i, j int) bool {
		return engines[i].GetCategory() < engines[j].GetCategory()
	})
}

func conjointTreat(text string) string {
	if strings.TrimSpace(text) == "" {
		return text
	}

	// 边缘空格去除
	text = strings.TrimSpace(text)
	if strings.TrimSpace(text) == "" {
		return text
	}

	// 首尾符号去除
	char := []rune(text)
	firstChar := char[0]
	for _const.CHAT_NOISE_COMMON_SET.Contains(string(firstChar)) {
		text = string(char[1:])
		char = []rune(text)
		firstChar = char[0]
	}

	if strings.TrimSpace(text) == "" {
		return text
	}

	lastChar := char[len(char)-1]
	for _const.CHAT_NOISE_COMMON_SET.Contains(string(lastChar)) {
		text = string(char[:len(char)-1])
		char = []rune(text)
		lastChar = char[len(char)-1]
	}
	return text
}
