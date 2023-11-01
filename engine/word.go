package engine

import (
	"slice-engine/entity"
	"slice-engine/enums"
)

type WordEngine struct {
	RootEngine
}

func init() {
	FACTORY = append(FACTORY, &WordEngine{})
}

func (e *WordEngine) GetCategory() enums.EngineType {
	return enums.EngineType_WORD
}

func (e *WordEngine) Recognize(state *entity.EngineState) bool {
	return false
}
