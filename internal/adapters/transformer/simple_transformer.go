package transformer

import (
	"context"
	"fmt"

	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	"github.com/xeniasokk/field-switcher/internal/ports"
	appErrors "github.com/xeniasokk/field-switcher/pkg/errors"
)

type SimpleTransformer struct {
	targetRole dream.Role
}

func NewSimpleTransformer(targetRole dream.Role) (*SimpleTransformer, error) {
	if targetRole == "" {
		return nil, appErrors.NewValidationError("target role cannot be empty")
	}
	return &SimpleTransformer{
		targetRole: targetRole,
	}, nil
}

var _ ports.TransformerPort = (*SimpleTransformer)(nil)

func (t *SimpleTransformer) TransformDream(
	ctx context.Context,
	child dream.ChildhoodDream,
) (dream.Adult, error) {
	_ = ctx

	switch child.Type() {
	case dream.TypeFootballer:
		return t.TransformFootballer(child)
	default:
		return dream.Adult{}, appErrors.NewDomainError(
			fmt.Sprintf("unsupported childhood dream type: %s", child.Type()),
		)
	}
}

func (t *SimpleTransformer) TransformFootballer(child dream.ChildhoodDream) (dream.Adult, error) {
	devField, err := dream.NewDevelopmentField()
	if err != nil {
		return dream.Adult{}, appErrors.Wrap(err, appErrors.CodeDomainFailure, "failed to create dev field")
	}

	roleConfig, err := t.targetRole.Config()
	if err != nil {
		return dream.Adult{}, appErrors.Wrap(err, appErrors.CodeDomainFailure, "failed to get role config")
	}

	stack := roleConfig.Stack()
	if len(stack) == 0 {
		return dream.Adult{}, appErrors.NewDomainError("role config must have a non-empty stack")
	}

	adult, err := dream.NewAdult(
		roleConfig.Title(),
		roleConfig.Description(),
		devField,
		stack,
		child.Qualities(),
		roleConfig.Comment(),
	)
	if err != nil {
		return dream.Adult{}, appErrors.Wrap(err, appErrors.CodeDomainFailure, "failed to create adult identity")
	}

	return adult, nil
}
