package boot

import (
	"reflect"
	"testing"
)

type args struct {
	val string
}

type test struct {
	name string
	args args
	want *[]string
}

func TestRun(t *testing.T) {
	tests := []test{
		{
			name: "ordinal",
			args: args{"1)、急性肝功能衰竭，(2)、急性胆囊炎；(3)、双肾结石"},
			want: &[]string{"急性肝功能衰竭", "急性胆囊炎", "双肾结石"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
