package shoot

import (
	"reflect"
	"testing"
)

func Test_noopFix(t *testing.T) {
	type args struct {
		src []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "oneline",
			args: args{
				src: []byte(`func Func() {}`),
			},
			want: []byte(`func Func() { /*noop*/ }`),
		},
		{
			name: "newline",
			args: args{
				src: []byte(`func Func() {
				}`),
			},
			want: []byte(`func Func() { /*noop*/ }`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noopFix(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("noopFix() = %v, want %v", got, tt.want)
			}
		})
	}
}
