package utils

import (
	mapset "github.com/deckarep/golang-set"
	"slice-engine/entity"
	"strings"
	"testing"
)

type args[T any] struct {
	set       mapset.Set
	predicate Predicate[T]
}
type testCase[T any] struct {
	name string
	args args[T]
	want bool
}

func before() []testCase[entity.PreSufTuple] {
	return []testCase[entity.PreSufTuple]{
		{
			name: "test1",
			args: args[entity.PreSufTuple]{
				set: mapset.NewSet(
					entity.NewPreSufTupleWithString("(", ")"),
					entity.NewPreSufTupleWithString("（", "）"),
					entity.NewPreSufTupleWithString("[", "]"),
				),
				predicate: func(tuple entity.PreSufTuple) bool {
					return strings.HasSuffix("", tuple.Pre) && strings.HasPrefix(")、急性肝功能衰竭，(2)、急性胆囊炎；(3)、双肾结石", tuple.Suf)
				},
			},
			want: false,
		},
	}
}

func TestAnySetMatch(t *testing.T) {
	for _, tt := range before() {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnySetMatch[entity.PreSufTuple](tt.args.set, tt.args.predicate); got != tt.want {
				t.Errorf("AnySetMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyMatch(t *testing.T) {
	for _, tt := range before() {
		slice := new([]entity.PreSufTuple)
		for v := range tt.args.set.Iter() {
			*slice = append(*slice, *v.(*entity.PreSufTuple))
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyMatch[entity.PreSufTuple](slice, tt.args.predicate); got != tt.want {
				t.Errorf("AnySetMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
