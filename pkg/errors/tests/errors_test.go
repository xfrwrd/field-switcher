package tests

import (
	stdErrors "errors"
	"testing"

	"github.com/xeniasokk/field-switcher/pkg/errors"
)

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    errors.Code
	}{
		{
			name:    "validation error",
			message: "test",
			want:    errors.CodeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.NewValidationError(tt.message)
			if err.Code() != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, err.Code())
			}
		})
	}
}

func TestWrapAndCodeOf(t *testing.T) {
	tests := []struct {
		name       string
		baseErr    error
		code       errors.Code
		message    string
		wantCode   errors.Code
		wantIsCode bool
	}{
		{
			name:       "wrap domain failure",
			baseErr:    stdErrors.New("base"),
			code:       errors.CodeDomainFailure,
			message:    "wrap",
			wantCode:   errors.CodeDomainFailure,
			wantIsCode: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.Wrap(tt.baseErr, tt.code, tt.message)
			if got := errors.CodeOf(err); got != tt.wantCode {
				t.Fatalf("expected CodeOf to return %v, got %v", tt.wantCode, got)
			}
			if got := errors.IsCode(err, tt.wantCode); got != tt.wantIsCode {
				t.Fatalf("expected IsCode to return %v, got %v", tt.wantIsCode, got)
			}
		})
	}
}
