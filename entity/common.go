package entity

import (
	"github.com/dlclark/regexp2"
	"regexp"
	"slice-engine/enums"
)

type Engine interface {
	//Run(state *EngineState)

	IsSupportFuzzy() bool

	GetCategory() enums.EngineType
	Recognize(state *EngineState) bool

	Slice(state *EngineState)
	Merge(state *EngineState)
	Optimize(state *EngineState)

	Standardize(state *EngineState)
	DeNoise(state *EngineState)
}

type OrdinalTuple struct {
	Text     string
	Seq      float64
	Suffix   string
	Position int
}

type OrdinalScore struct {
	Score     float64
	IsPerfect bool
}

type PatternReplaceTuple struct {
	Matter  *regexp2.Regexp
	Replace string
}

type PreSufTuple struct {
	Pre       string
	Suf       string
	PreRegexp *regexp.Regexp
	SufRegexp *regexp.Regexp
}

type SemicolonTuple struct {
	Text string
	Sem  string
}

func NewPreSufTupleWithString(pre, suf string) *PreSufTuple {
	return &PreSufTuple{
		Pre:       pre,
		Suf:       suf,
		PreRegexp: nil,
		SufRegexp: nil,
	}
}

func NewPreSufTupleWithRegexp(preRegexp, sufRegexp *regexp.Regexp) *PreSufTuple {
	return &PreSufTuple{
		Pre:       "",
		Suf:       "",
		PreRegexp: preRegexp,
		SufRegexp: sufRegexp,
	}
}
