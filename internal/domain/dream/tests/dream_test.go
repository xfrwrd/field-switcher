package tests

import (
	"testing"

	"github.com/xeniasokk/field-switcher/internal/domain/dream"
)

func TestNewQualityValidation(t *testing.T) {
	tests := []struct {
		name        string
		qualityName string
		description string
		wantErr     bool
	}{
		{
			name:        "empty name should fail",
			qualityName: "",
			description: "desc",
			wantErr:     true,
		},
		{
			name:        "valid quality",
			qualityName: "Quality",
			description: "Description",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dream.NewQuality(tt.qualityName, tt.description)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewQuality() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAdultCopiesSlices(t *testing.T) {
	f, _ := dream.NewField("Field", "Env")
	q, _ := dream.NewQuality("Q", "D")
	stack := []string{"Go"}
	traits := []dream.Quality{q}

	a, err := dream.NewAdult("Role", "Desc", f, stack, traits, "comment")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Modify original slices
	stack[0] = "JS"
	traits[0] = dream.Quality{}

	tests := []struct {
		name     string
		got      interface{}
		want     interface{}
		checkFn  func() bool
		errorMsg string
	}{
		{
			name:     "stack should be copied",
			got:      a.Stack()[0],
			want:     "Go",
			checkFn:  func() bool { return a.Stack()[0] == "Go" },
			errorMsg: "expected stack copy",
		},
		{
			name:     "traits should be copied",
			got:      a.Traits()[0].Name(),
			want:     "Q",
			checkFn:  func() bool { return a.Traits()[0].Name() == "Q" },
			errorMsg: "expected traits copy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.checkFn() {
				t.Fatalf("%s: got %v, want %v", tt.errorMsg, tt.got, tt.want)
			}
		})
	}
}

func TestRoleConfig(t *testing.T) {
	tests := []struct {
		name        string
		role        dream.Role
		wantTitle   string
		wantErr     bool
		checkFields func(dream.RoleConfig) bool
	}{
		{
			name:      "team lead role",
			role:      dream.RoleTeamLead,
			wantTitle: "Тимлид",
			wantErr:   false,
			checkFields: func(cfg dream.RoleConfig) bool {
				return cfg.Description() != "" && cfg.Comment() != "" && len(cfg.Stack()) > 0
			},
		},
		{
			name:      "developer role",
			role:      dream.RoleDeveloper,
			wantTitle: "Разработчик",
			wantErr:   false,
			checkFields: func(cfg dream.RoleConfig) bool {
				return cfg.Description() != "" && cfg.Comment() != "" && len(cfg.Stack()) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tt.role.Config()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Config() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if config.Title() != tt.wantTitle {
					t.Fatalf("expected title '%s', got '%s'", tt.wantTitle, config.Title())
				}
				if !tt.checkFields(config) {
					t.Fatalf("expected non-empty description and comment")
				}
			}
		})
	}
}

func TestRoleConfigUnsupported(t *testing.T) {
	tests := []struct {
		name    string
		role    dream.Role
		wantErr bool
	}{
		{
			name:    "unsupported role",
			role:    dream.Role("unsupported"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.role.Config()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Config() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewDevelopmentField(t *testing.T) {
	tests := []struct {
		name            string
		wantName        string
		wantEnvironment string
		wantErr         bool
	}{
		{
			name:            "development field",
			wantName:        dream.DevFieldName,
			wantEnvironment: dream.DevFieldEnvironment,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, err := dream.NewDevelopmentField()
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewDevelopmentField() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if field.Name() != tt.wantName {
					t.Fatalf("expected field name '%s', got '%s'", tt.wantName, field.Name())
				}
				if field.Environment() != tt.wantEnvironment {
					t.Fatalf("expected field environment '%s', got '%s'", tt.wantEnvironment, field.Environment())
				}
			}
		})
	}
}
