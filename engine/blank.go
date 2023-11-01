package engine

import (
	mapset "github.com/deckarep/golang-set"
	"regexp"
	_const "slice-engine/const"
	"slice-engine/entity"
	"slice-engine/enums"
	"slice-engine/utils"
	"strings"
)

var (
	SPLITTER_PATTERN      = regexp.MustCompile("(\\s+)")
	CHAT_JOIN_SET         = mapset.NewSet("-", "—", "－", "―")
	UNIT_JOIN_PATTERN_STR = "(" + strings.Join(utils.Transform[string](CHAT_JOIN_SET.ToSlice()), "|") + ")"

	ILLEGAL_PATTERN_TUPLE_SET = mapset.NewSet(
		entity.PreSufTuple{
			PreRegexp: regexp.MustCompile("[\\-\\_\\+第孕产]$"),
			SufRegexp: regexp.MustCompile("^\\d"),
		},
		entity.PreSufTuple{
			PreRegexp: regexp.MustCompile("\\d$"),
			SufRegexp: regexp.MustCompile("^[" + strings.Join(utils.Transform[string](_const.UNIT_SUFFIX_SET.ToSlice()), "") + "]"),
		},
	)
)

type BlankEngine struct {
	RootEngine
}

func init() {
	FACTORY = append(FACTORY, &BlankEngine{})
}

func (e *BlankEngine) GetCategory() enums.EngineType {
	return enums.EngineType_BLANK
}

func (e *BlankEngine) Recognize(state *entity.EngineState) bool {
	text := state.Text
	if text == "" || len(text) == 0 {
		return false
	}

	return false
}
