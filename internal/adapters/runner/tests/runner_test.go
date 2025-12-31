package tests

import (
	"bytes"
	"context"
	"testing"

	"github.com/xeniasokk/field-switcher/internal/adapters/formatter"
	"github.com/xeniasokk/field-switcher/internal/adapters/presenter"
	"github.com/xeniasokk/field-switcher/internal/adapters/runner"
	"github.com/xeniasokk/field-switcher/internal/adapters/transformer"
	"github.com/xeniasokk/field-switcher/internal/application/usecase/transform"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
)

func TestConsoleRunnerRun(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		targetRole dream.Role
		wantOutput bool
		wantErr    bool
	}{
		{
			name:       "run with developer role",
			targetRole: dream.RoleDeveloper,
			wantOutput: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := transformer.NewSimpleTransformer(tt.targetRole)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			uc := transform.NewUseCase(tr)
			p := presenter.NewConsolePresenter()
			f := formatter.NewTextFormatter()

			var buf bytes.Buffer

			r, err := runner.NewConsoleRunner(uc, p, f, &buf)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			child, err := dream.NewDefaultFootballerDream()
			if err != nil {
				t.Fatalf("failed to create default footballer dream: %v", err)
			}

			err = r.Run(ctx, child)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.wantOutput {
				if buf.Len() == 0 {
					t.Fatalf("expected some output")
				}
			}
		})
	}
}
