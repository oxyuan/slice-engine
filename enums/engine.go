package enums

import (
	mapset "github.com/deckarep/golang-set"
	"math"
)

type EngineType int

const (
	EngineType_ORDINAL EngineType = iota
	EngineType_CODE
	EngineType_SEMICOLOW
	EngineType_BLANK
	EngineType_WORD
	EngineType_NONE
)

type EG struct {
	Priority int
	Excludes mapset.Set
}

var EngineEnum = map[EngineType]EG{
	EngineType_ORDINAL:   {1, mapset.NewSet()},
	EngineType_CODE:      {2, mapset.NewSet()},
	EngineType_SEMICOLOW: {3, mapset.NewSet()},
	EngineType_BLANK:     {4, mapset.NewSet()},
	EngineType_WORD:      {5, mapset.NewSet(DataType_OPT)},
	EngineType_NONE:      {math.MaxInt32, mapset.NewSet()},
}
