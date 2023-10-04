package splitter

import (
	mapset "github.com/deckarep/golang-set"
	"reflect"
	"regexp"
	"slice-engine/entity"
	"testing"
)

func TestRegex(t *testing.T) {
	illegalPatternTupleSet = mapset.NewSet(
		entity.NewPreSufTupleWithRegexp(
			regexp.MustCompile("[\\(\\[（【][\\w\\.；]*$"),
			regexp.MustCompile("^[\\w\\.;；]*[\\)\\]）】]")),
	)

	for val := range illegalPatternTupleSet.Iter() {
		tuple := val.(*entity.PreSufTuple)
		t.Log(tuple.PreRegexp)
		t.Log(tuple.SufRegexp)
	}
}

func TestSplit(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []*entity.OrdinalTuple
	}{
		{
			name: "split.func()",
			args: args{input: "1)、急性肝功能衰竭，(2)、急性胆囊炎；(3)、双肾结石"},
			want: []*entity.OrdinalTuple{
				{
					Text:     "1)、急性肝功能衰竭，(",
					Seq:      1,
					Suffix:   ")",
					Position: 2,
				},
				{
					Text:     "2)、急性胆囊炎；(",
					Seq:      2,
					Suffix:   ")",
					Position: 2,
				},
				{
					Text:     "3)、双肾结石",
					Seq:      3,
					Suffix:   ")",
					Position: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Split(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assignTuple(t *testing.T) {
	type args struct {
		list *[]*entity.OrdinalTuple
		text string
		seq  string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assignTuple(tt.args.list, tt.args.text, tt.args.seq)
		})
	}
}

func Test_getSeq(t *testing.T) {
	type args struct {
		seq string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSeq(tt.args.seq); got != tt.want {
				t.Errorf("getSeq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasLetter(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasLetter(tt.args.prefix, tt.args.suffix); got != tt.want {
				t.Errorf("hasLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isAllowed(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAllowed(tt.args.prefix, tt.args.suffix); got != tt.want {
				t.Errorf("isAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isIllegal(t *testing.T) {
	type args struct {
		prefix string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIllegal(tt.args.prefix, tt.args.suffix); got != tt.want {
				t.Errorf("isIllegal() = %v, want %v", got, tt.want)
			}
		})
	}
}
