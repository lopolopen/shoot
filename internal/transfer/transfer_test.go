package transfer_test

import (
	"testing"

	"github.com/lopolopen/shoot/internal/transfer"
)

func TestFirstLower(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: "",
		},
		{
			name: "lowercase",
			s:    "lowercase",
			want: "l",
		},
		{
			name: "Uppercase",
			s:    "Uppercase",
			want: "u",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transfer.FirstLower(tt.s)
			if got != tt.want {
				t.Errorf("FirstLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: "",
		},
		{
			name: "camel case",
			s:    "camelCase",
			want: "camelCase",
		},
		{
			name: "pascal case",
			s:    "PascalCase",
			want: "pascalCase",
		},
		{
			name: "snake case",
			s:    "snake_case",
			want: "snakeCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transfer.ToCamelCase(tt.s)
			if got != tt.want {
				t.Errorf("ToCamelCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: "",
		},
		{
			name: "pascal case",
			s:    "PascalCase",
			want: "PascalCase",
		},
		{
			name: "camel case",
			s:    "camelCase",
			want: "CamelCase",
		},
		{
			name: "snake case",
			s:    "snake_case",
			want: "SnakeCase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := transfer.ToPascalCase(tt.s)
			if got != tt.want {
				t.Errorf("ToPascalCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
