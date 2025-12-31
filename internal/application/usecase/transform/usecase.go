package transform

import (
	"context"

	"github.com/xeniasokk/field-switcher/internal/ports"
	"github.com/xeniasokk/field-switcher/pkg/errors"
)

type UseCase interface {
	Execute(ctx context.Context, input InputModel) (OutputModel, error)
}

type useCase struct {
	transformer ports.TransformerPort
}

func NewUseCase(transformer ports.TransformerPort) UseCase {
	return &useCase{
		transformer: transformer,
	}
}

func (uc *useCase) Execute(ctx context.Context, input InputModel) (OutputModel, error) {
	if err := input.Validate(ctx); err != nil {
		return OutputModel{}, errors.Wrap(err, errors.CodeValidation, "invalid input for TransformDreamUseCase")
	}

	adult, err := uc.transformer.TransformDream(ctx, input.Child())
	if err != nil {
		return OutputModel{}, errors.Wrap(err, errors.CodeDomainFailure, "transformer failed to process dream")
	}

	out := NewOutput(input.Child(), adult)
	return out, nil
}
