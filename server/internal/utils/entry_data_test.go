package utils_test

import (
	"testing"
	"time"

	"github.com/boreymarf/task-fuss/server/internal/utils"
)

func TestEvaluateCondition(t *testing.T) {
	tests := []struct {
		name      string
		condition string
		args      []any
		want      bool
		wantErr   bool
	}{
		{
			name:      "== integers true",
			condition: "==",
			args:      []any{10, 10, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "== integers false",
			condition: "==",
			args:      []any{10, 20, 30},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "== strings true",
			condition: "==",
			args:      []any{"hello", "hello", "hello"},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "== strings false",
			condition: "==",
			args:      []any{"hello", "world", "test"},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "== types false",
			condition: "==",
			args:      []any{10, "10", 10.0},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "== insufficient arguments",
			condition: "==",
			args:      []any{10},
			want:      false,
			wantErr:   true,
		},
		{
			name:      "== empty arguments",
			condition: "==",
			args:      []any{},
			want:      false,
			wantErr:   true,
		},

		{
			name:      ">= int true equal",
			condition: ">=",
			args:      []any{10, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      ">= int true greater",
			condition: ">=",
			args:      []any{20, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      ">= int false",
			condition: ">=",
			args:      []any{5, 10},
			want:      false,
			wantErr:   false,
		},
		{
			name:      ">= duration true",
			condition: ">=",
			args:      []any{time.Minute, 30 * time.Second},
			want:      true,
			wantErr:   false,
		},
		{
			name:      ">= duration false",
			condition: ">=",
			args:      []any{time.Second, time.Minute},
			want:      false,
			wantErr:   false,
		},
		{
			name:      ">= wrong argument count",
			condition: ">=",
			args:      []any{10},
			want:      false,
			wantErr:   true,
		},
		{
			name:      ">= too many arguments",
			condition: ">=",
			args:      []any{10, 20, 30},
			want:      false,
			wantErr:   true,
		},
		{
			name:      ">= mixed types",
			condition: ">=",
			args:      []any{10, "10"},
			want:      false,
			wantErr:   true,
		},
		{
			name:      ">= unsupported type",
			condition: ">=",
			args:      []any{"a", "b"},
			want:      false,
			wantErr:   true,
		},

		{
			name:      "> int true",
			condition: ">",
			args:      []any{20, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "> int false equal",
			condition: ">",
			args:      []any{10, 10},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "> int false less",
			condition: ">",
			args:      []any{5, 10},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "> duration true",
			condition: ">",
			args:      []any{time.Minute, 30 * time.Second},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "> duration false",
			condition: ">",
			args:      []any{time.Second, time.Minute},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "> wrong argument count",
			condition: ">",
			args:      []any{10},
			want:      false,
			wantErr:   true,
		},
		{
			name:      "> mixed types",
			condition: ">",
			args:      []any{10, "10"},
			want:      false,
			wantErr:   true,
		},
		{
			name:      "> unsupported type",
			condition: ">",
			args:      []any{"a", "b"},
			want:      false,
			wantErr:   true,
		},

		{
			name:      "< int true",
			condition: "<",
			args:      []any{5, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "< int false equal",
			condition: "<",
			args:      []any{10, 10},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "< int false greater",
			condition: "<",
			args:      []any{20, 10},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "< duration true",
			condition: "<",
			args:      []any{time.Second, time.Minute},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "< duration false",
			condition: "<",
			args:      []any{time.Minute, 30 * time.Second},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "< wrong argument count",
			condition: "<",
			args:      []any{10},
			want:      false,
			wantErr:   true,
		},
		{
			name:      "< mixed types",
			condition: "<",
			args:      []any{10, "10"},
			want:      false,
			wantErr:   true,
		},

		{
			name:      "<= int true less",
			condition: "<=",
			args:      []any{5, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "<= int true equal",
			condition: "<=",
			args:      []any{10, 10},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "<= int false",
			condition: "<=",
			args:      []any{20, 10},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "<= duration true",
			condition: "<=",
			args:      []any{time.Second, time.Minute},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "<= duration true equal",
			condition: "<=",
			args:      []any{time.Minute, time.Minute},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "<= duration false",
			condition: "<=",
			args:      []any{2 * time.Minute, time.Minute},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "<= wrong argument count",
			condition: "<=",
			args:      []any{10},
			want:      false,
			wantErr:   true,
		},
		{
			name:      "<= mixed types",
			condition: "<=",
			args:      []any{10, "10"},
			want:      false,
			wantErr:   true,
		},

		{
			name:      "AND all true",
			condition: "AND",
			args:      []any{true, true, true},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "AND one false",
			condition: "AND",
			args:      []any{true, false, true},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "OR all false",
			condition: "OR",
			args:      []any{false, false, false},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "OR one true",
			condition: "OR",
			args:      []any{false, true, false},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "NAND all true",
			condition: "NAND",
			args:      []any{true, true, true},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "NAND one false",
			condition: "NAND",
			args:      []any{true, false, true},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "NOR all false",
			condition: "NOR",
			args:      []any{false, false, false},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "NOR one true",
			condition: "NOR",
			args:      []any{false, true, false},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "single argument AND",
			condition: "AND",
			args:      []any{true},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "single argument OR",
			condition: "OR",
			args:      []any{false},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "no arguments logical",
			condition: "AND",
			args:      []any{},
			want:      false,
			wantErr:   true,
		},
		{
			name:      "non-bool argument",
			condition: "AND",
			args:      []any{true, "not bool"},
			want:      false,
			wantErr:   true,
		},

		{
			name:      "unsupported condition",
			condition: "nil",
			args:      []any{10, 20},
			want:      false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.EvaluateCondition(tt.condition, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvaluateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvaluateCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
