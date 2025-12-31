package tests

import (
	"context"
	"testing"

	"github.com/xeniasokk/field-switcher/internal/adapters/transformer"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
)

func TestNewSimpleTransformerValidation(t *testing.T) {
	tests := []struct {
		name       string
		targetRole dream.Role
		wantErr    bool
	}{
		{
			name:       "empty target role should fail",
			targetRole: dream.Role(""),
			wantErr:    true,
		},
		{
			name:       "valid transformer",
			targetRole: dream.RoleDeveloper,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := transformer.NewSimpleTransformer(tt.targetRole)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewSimpleTransformer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSimpleTransformerTransformDream(t *testing.T) {
	ctx := context.Background()
	child, err := dream.NewDefaultFootballerDream()
	if err != nil {
		t.Fatalf("failed to create default footballer dream: %v", err)
	}

	tests := []struct {
		name         string
		targetRole   dream.Role
		wantTitle    string
		wantStackLen int
		wantErr      bool
		checkFn      func(dream.Adult) bool
	}{
		{
			name:         "transform to developer",
			targetRole:   dream.RoleDeveloper,
			wantTitle:    "Разработчик",
			wantStackLen: 3, // Стек из роли: ["Go", "Git", "Microservices"]
			wantErr:      false,
			checkFn: func(a dream.Adult) bool {
				stack := a.Stack()
				return a.RoleTitle() != "" && len(stack) == 3 && stack[0] == "Go"
			},
		},
		{
			name:         "transform to team lead",
			targetRole:   dream.RoleTeamLead,
			wantTitle:    "Тимлид",
			wantStackLen: 7, // Стек из роли: 7 элементов для тимлида
			wantErr:      false,
			checkFn: func(a dream.Adult) bool {
				stack := a.Stack()
				return a.RoleTitle() == "Тимлид" && len(stack) == 7 && stack[0] == "System Design" && a.Comment() != ""
			},
		},
		{
			name:         "unsupported role should fail",
			targetRole:   dream.Role("unknown_role"),
			wantTitle:    "",
			wantStackLen: 0,
			wantErr:      true,
			checkFn:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := transformer.NewSimpleTransformer(tt.targetRole)
			if err != nil {
				t.Fatalf("unexpected error creating transformer: %v", err)
			}

			adult, err := tr.TransformDream(ctx, child)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TransformDream() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if adult.RoleTitle() != tt.wantTitle {
					t.Fatalf("expected role title '%s', got '%s'", tt.wantTitle, adult.RoleTitle())
				}
				if len(adult.Stack()) != tt.wantStackLen {
					t.Fatalf("expected stack length %d, got %d", tt.wantStackLen, len(adult.Stack()))
				}
				if tt.checkFn != nil && !tt.checkFn(adult) {
					t.Fatalf("check function failed")
				}
			}
		})
	}
}
