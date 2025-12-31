package transform

import (
	"context"

	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	domainErrors "github.com/xeniasokk/field-switcher/pkg/errors"
)

type InputModel struct {
	child dream.ChildhoodDream
}

func NewInput(child dream.ChildhoodDream) InputModel {
	return InputModel{child: child}
}

func (i InputModel) Child() dream.ChildhoodDream {
	return i.child
}

func (i InputModel) Validate(ctx context.Context) error {
	_ = ctx
	if i.child.DisplayName() == "" {
		return domainErrors.NewValidationError("childhood dream has no display name")
	}
	return nil
}
