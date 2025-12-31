package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/xeniasokk/field-switcher/internal/adapters/formatter"
	"github.com/xeniasokk/field-switcher/internal/adapters/presenter"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
)

func TestTextFormatterFormat(t *testing.T) {
	ctx := context.Background()
	child, err := dream.NewDefaultFootballerDream()
	if err != nil {
		t.Fatalf("failed to create default footballer dream: %v", err)
	}
	f, err := dream.NewField("Dev", "Team")
	if err != nil {
		t.Fatalf("failed to create field: %v", err)
	}
	a, err := dream.NewAdult("Role", "Desc", f, []string{"Go"}, child.Qualities(), "comment")
	if err != nil {
		t.Fatalf("failed to create adult: %v", err)
	}

	tests := []struct {
		name         string
		vm           func() presenter.ConsoleViewModel
		wantContains []string
		wantErr      bool
	}{
		{
			name: "format with all sections",
			vm: func() presenter.ConsoleViewModel {
				return presenter.NewConsoleViewModel("title", child, a, "note")
			},
			wantContains: []string{"ДЕТСКАЯ МЕЧТА", "ВЗРОСЛАЯ РОЛЬ"},
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := formatter.NewTextFormatter()
			text, err := formatter.Format(ctx, tt.vm())
			if (err != nil) != tt.wantErr {
				t.Fatalf("Format() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				for _, want := range tt.wantContains {
					if !strings.Contains(text, want) {
						t.Fatalf("expected text to contain '%s'", want)
					}
				}
			}
		})
	}
}
