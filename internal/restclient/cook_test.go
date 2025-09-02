package restclient

import (
	"reflect"
	"testing"
)

func Test_parsePath(t *testing.T) {
	type args struct {
		doc string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 []string
		want3 bool
	}{
		{
			name: "normal case with quotes",
			args: args{
				doc: `shoot: Get("/api/users{id}")`,
			},
			want:  "GET",
			want1: "/api/users{id}",
			want2: []string{"id"},
			want3: true,
		},
		{
			name: "normal case without quotes",
			args: args{
				doc: `shoot: Get(/api/users{id})`,
			},
			want:  "GET",
			want1: "/api/users{id}",
			want2: []string{"id"},
			want3: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := parsePath(tt.args.doc)
			if got != tt.want {
				t.Errorf("parsePath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parsePath() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("parsePath() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("parsePath() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}
