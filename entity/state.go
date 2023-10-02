package entity

import (
	_const "slice-engine/const"
	"slice-engine/enums"
)

type EngineState struct {
	Text           string
	Type           *enums.EngineType
	IsRecognized   bool
	Group          *[]string
	MergeList      *[]*[]string
	IsFuzzy        bool
	Chain          *[]*Engine
	IsUntapped     bool
	SplitList      *[]string
	SemicolonTuple *[]*SemicolonTuple
	OrdinalTuple   *[]*OrdinalTuple
}

/*func (e *EngineState) AddChain(ee *EG) {
	if e.Chain == nil {
		e.Chain = new([]*EG)
	}
	// 如果 engine 是 NoneEngine 类型，不执行 append 操作
	if _, ok := (*ee).(*engine.NoneEngine); ok {
		return
	}
	if reflect.TypeOf(ee).String() == "*ee.NoneEngine" {
		return
	}
	*e.Chain = append(*e.Chain, ee)
}*/

func (e *EngineState) IsBreak() bool {
	if e.IsUntapped {
		return true
	}
	size := len(*e.Chain)
	if size > _const.DEF_MAX_CHAIN_DEPTH {
		return true
	}
	if size > 3 {
		one := (*e.Chain)[0]
		two := (*e.Chain)[1]
		three := (*e.Chain)[2]
		return one == two && two == three
	}
	return false
}

func (e *EngineState) Clear() {
	//e.Text = nil
	//e.Type = nil
	//e.IsRecognized = nil
	e.Group = new([]string)
	e.MergeList = new([]*[]string)
	//e.IsFuzzy = nil
	e.Chain = new([]*Engine)
	//e.IsUntapped = nil
	e.SplitList = new([]string)
	e.SemicolonTuple = new([]*SemicolonTuple)
	e.OrdinalTuple = new([]*OrdinalTuple)
}

func (e *EngineState) IsExcluded(ee *Engine) bool {
	category := (*ee).GetCategory()
	eg := enums.EngineEnum[category]
	if eg.Excludes == nil || len(eg.Excludes.ToSlice()) == 0 {
		return false
	}
	return eg.Excludes.Contains(e.Type)
}

func NewEngineState(text string) *EngineState {
	return &EngineState{
		Text:         text,
		Type:         nil,
		IsRecognized: false,
		IsUntapped:   true,
		IsFuzzy:      false,

		Group:          new([]string),
		MergeList:      new([]*[]string),
		Chain:          new([]*Engine),
		SplitList:      new([]string),
		SemicolonTuple: new([]*SemicolonTuple),
		OrdinalTuple:   new([]*OrdinalTuple),
	}
}
