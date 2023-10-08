package engine

import (
	"slice-engine/entity"
	"slice-engine/enums"
)

type SemicolonEngine struct {
	RootEngine
}

func init() {
	FACTORY = append(FACTORY, &SemicolonEngine{})
}

func (e *SemicolonEngine) GetCategory() enums.EngineType {
	return enums.EngineType_SEMICOLOW
}

// Recognize 识别标点分位点
//
// 1. 获取出现频次最多的分位点做分割；
// 2. 在括号中出现 有且只有一个 分位点，此位点不能分割；
// 3. 位点后不能出现等级词；
// 4. 两个位点间不能只有一个字词；
func (e *SemicolonEngine) Recognize(state *entity.EngineState) bool {
	text := state.Text
	if text == "" {
		return false
	}

	return false
}
