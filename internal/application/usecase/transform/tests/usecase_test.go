package tests

import (
	"context"
	"testing"

	"github.com/xeniasokk/field-switcher/internal/application/usecase/transform"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	"github.com/xeniasokk/field-switcher/internal/ports"
)

type transformerMock struct {
	adult dream.Adult
	err   error
}

func (m *transformerMock) TransformDream(ctx context.Context, d dream.ChildhoodDream) (dream.Adult, error) {
	_ = ctx
	return m.adult, m.err
}

var _ ports.TransformerPort = (*transformerMock)(nil)

func TestUseCaseExecute(t *testing.T) {
	ctx := context.Background()
	child, err := dream.NewDefaultFootballerDream()
	if err != nil {
		t.Fatalf("failed to create default footballer dream: %v", err)
	}

	tests := []struct {
		name      string
		mockAdult dream.Adult
		mockErr   error
		wantTitle string
		wantErr   bool
	}{
		{
			name: "successful execution",
			mockAdult: func() dream.Adult {
				f, err := dream.NewField("Dev", "Team")
				if err != nil {
					t.Fatalf("failed to create field: %v", err)
				}
				a, err := dream.NewAdult("Role", "Desc", f, nil, nil, "")
				if err != nil {
					t.Fatalf("failed to create adult: %v", err)
				}
				return a
			}(),
			mockErr:   nil,
			wantTitle: "Role",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &transformerMock{
				adult: tt.mockAdult,
				err:   tt.mockErr,
			}
			uc := transform.NewUseCase(mock)

			input := transform.NewInput(child)

			out, err := uc.Execute(ctx, input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if out.Adult().RoleTitle() != tt.wantTitle {
					t.Fatalf("expected adult role '%s', got '%s'", tt.wantTitle, out.Adult().RoleTitle())
				}
			}
		})
	}
}
