package engine

import (
	"slice-engine/entity"
	"slice-engine/enums"
)

type NoneEngine struct {
	RootEngine
}

func init() {
	FACTORY = append(FACTORY, &NoneEngine{})
}

func (e *NoneEngine) GetCategory() enums.EngineType {
	return enums.EngineType_NONE
}

func (e *NoneEngine) Recognize(state *entity.EngineState) bool {
	return false
}

func (e *NoneEngine) IsSupportFuzzy() bool {
	return true
}
