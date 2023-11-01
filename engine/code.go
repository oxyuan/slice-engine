package engine

import (
	"regexp"
	"slice-engine/entity"
	"slice-engine/enums"
	"slice-engine/utils"
)

var (
	splitterPattern = regexp.MustCompile("(\\([\\w\\.]{5,}\\))")
)

type CodeEngine struct {
	RootEngine
}

func init() {
	FACTORY = append(FACTORY, &CodeEngine{})
}

func (e *CodeEngine) GetCategory() enums.EngineType {
	return enums.EngineType_CODE
}

func (e *CodeEngine) Recognize(state *entity.EngineState) bool {
	text := state.Text
	if text == "" {
		return false
	}
	return utils.CountMatches(text, splitterPattern) > 1
}

func (e *CodeEngine) Slice(state *entity.EngineState) {
	text := state.Text
	if text != "" {
		state.Group = utils.Slice(text, splitterPattern)
	}
	return
}

func (e *CodeEngine) Merge(state *entity.EngineState) {
	list := make([]*[]string, 0)
	list = append(list, state.Group)
	state.MergeList = &list
}
