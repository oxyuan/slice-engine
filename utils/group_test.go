package utils

import (
	"regexp"
	"testing"
)

func TestSlice(t *testing.T) {
	type args struct {
		text   string
		regexp *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want *[]string
	}{
		{
			name: "test1",
			args: args{
				text:   "(060.000x001)G3P132+6周孕先兆早产(042.900)胎膜早破(013.x00x001)妊娠期高血压(034.200x002)剖宫产史的妊娠",
				regexp: regexp.MustCompile("(\\([\\w\\.]{5,}\\))"),
			},
			want: &[]string{"(060.000x001)G3P132+6周孕先兆早产", "(042.900)胎膜早破", "(013.x00x001)妊娠期高血压", "(034.200x002)剖宫产史的妊娠"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slice(tt.args.text, tt.args.regexp); !AreSlicesEqual(got, tt.want) {
				t.Errorf("Slice() = %v, want %v", got, tt.want)
			}
		})
	}
}
