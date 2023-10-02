package engine

import (
	"reflect"
	"slice-engine/entity"
)

var FACTORY = make([]entity.Engine, 0)

type RootEngine struct {
	entity.Engine
}

func (e *RootEngine) Run(state *entity.EngineState) {
	// 二轮识别
	if e.isStateSupportFuzzy(state) {
		state.IsRecognized = false
		return
	}
	// 识别
	recognize := e.Engine.Recognize(state)
	state.IsRecognized = recognize
	if !recognize {
		state.Clear()
		return
	}

	// 分组
	e.group(state)

	if !state.IsRecognized {
		state.Clear()
		return
	}
	addChain(state, &e.Engine)
	if _, ok := e.Engine.(*NoneEngine); ok {
		state.IsUntapped = true
	} else {
		state.IsUntapped = false
	}
}

func (e *RootEngine) close() {

}

func (e *RootEngine) group(state *entity.EngineState) {
	// 1. slice the text into groups using the rule.
	e.Engine.Slice(state)
	// 2. merge groups that belong to the same serial number.
	e.Engine.Merge(state)
	// 3. select the group with the highest compliance.
	e.Engine.Optimize(state)
}

func (e *RootEngine) IsSupportFuzzy() bool {
	return false
}

func (e *RootEngine) isStateSupportFuzzy(state *entity.EngineState) bool {
	return state.IsFuzzy && !e.Engine.IsSupportFuzzy()
}

func (e *RootEngine) Slice(state *entity.EngineState) {
	state.Group = &[]string{state.Text}
}

func (e *RootEngine) Merge(state *entity.EngineState) {
	state.MergeList = &[]*[]string{state.Group}
}

func (e *RootEngine) Optimize(state *entity.EngineState) {
	state.Group = (*state.MergeList)[0]
}

func (e *RootEngine) Standardize(state *entity.EngineState) {

}
func (e *RootEngine) DeNoise(state *entity.EngineState) {

}

func addChain(es *entity.EngineState, ee *entity.Engine) {
	if es.Chain == nil {
		es.Chain = new([]*entity.Engine)
	}
	// 如果 engine 是 NoneEngine 类型，不执行 append 操作
	if _, ok := (*ee).(*NoneEngine); ok {
		return
	}
	if reflect.TypeOf(ee).String() == "*ee.NoneEngine" {
		return
	}
	*es.Chain = append(*es.Chain, ee)
}
