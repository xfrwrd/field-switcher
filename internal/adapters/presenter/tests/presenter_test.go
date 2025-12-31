package tests

import (
	"context"
	"testing"

	"github.com/xeniasokk/field-switcher/internal/adapters/presenter"
	"github.com/xeniasokk/field-switcher/internal/application/usecase/transform"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
)

func TestConsolePresenterPresent(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		child     dream.ChildhoodDream
		adult     dream.Adult
		wantTitle string
		wantErr   bool
	}{
		{
			name: "present with valid data",
			child: func() dream.ChildhoodDream {
				child, err := dream.NewDefaultFootballerDream()
				if err != nil {
					t.Fatalf("failed to create default footballer dream: %v", err)
				}
				return child
			}(),
			adult: func() dream.Adult {
				f, err := dream.NewField("Dev", "Team")
				if err != nil {
					t.Fatalf("failed to create field: %v", err)
				}
				child, err := dream.NewDefaultFootballerDream()
				if err != nil {
					t.Fatalf("failed to create default footballer dream: %v", err)
				}
				a, err := dream.NewAdult("Role", "Desc", f, nil, child.Qualities(), "comment")
				if err != nil {
					t.Fatalf("failed to create adult: %v", err)
				}
				return a
			}(),
			wantTitle: "Role",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := transform.NewOutput(tt.child, tt.adult)

			p := presenter.NewConsolePresenter()
			vmRaw, err := p.Present(ctx, out)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Present() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				vm, ok := vmRaw.(presenter.ConsoleViewModel)
				if !ok {
					t.Fatalf("expected ConsoleViewModel, got %T", vmRaw)
				}
				if vm.Adult().RoleTitle() != tt.wantTitle {
					t.Fatalf("expected role title '%s', got '%s'", tt.wantTitle, vm.Adult().RoleTitle())
				}
			}
		})
	}
}
