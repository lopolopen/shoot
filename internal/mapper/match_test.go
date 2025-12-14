package mapper

import (
	"go/types"
	"testing"
)

func Test_matchType(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		type1 types.Type
		type2 types.Type
		want  bool
		want2 bool
	}{
		{
			name:  "same1",
			type1: types.Typ[types.Int],
			type2: types.Typ[types.Int],
			want:  true,
			want2: true,
		},
		{
			name:  "same2",
			type1: types.Typ[types.Int],
			type2: types.Typ[types.Int32],
			want:  false,
			want2: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got2 := matchType(tt.type1, tt.type2)
			if got != tt.want {
				t.Errorf("matchType() = %v, want %v", got, tt.want)
			}
			if got2 != tt.want2 {
				t.Errorf("matchType() = %v, want %v", got2, tt.want2)
			}
		})
	}
}
