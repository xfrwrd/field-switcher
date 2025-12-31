package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/xeniasokk/field-switcher/internal/application/usecase/transform"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	"github.com/xeniasokk/field-switcher/internal/ports"
	appErrors "github.com/xeniasokk/field-switcher/pkg/errors"
)

type ConsoleRunner struct {
	useCase   transform.UseCase
	presenter ports.PresenterPort
	formatter ports.FormatterPort
	out       io.Writer
}

func NewConsoleRunner(
	useCase transform.UseCase,
	presenter ports.PresenterPort,
	formatter ports.FormatterPort,
	out io.Writer,
) (*ConsoleRunner, error) {
	if useCase == nil {
		return nil, appErrors.NewInternalError("useCase cannot be nil in ConsoleRunner")
	}
	if presenter == nil {
		return nil, appErrors.NewInternalError("presenter cannot be nil in ConsoleRunner")
	}
	if formatter == nil {
		return nil, appErrors.NewInternalError("formatter cannot be nil in ConsoleRunner")
	}
	if out == nil {
		return nil, appErrors.NewInternalError("output writer cannot be nil in ConsoleRunner")
	}

	return &ConsoleRunner{
		useCase:   useCase,
		presenter: presenter,
		formatter: formatter,
		out:       out,
	}, nil
}

func (r *ConsoleRunner) Run(ctx context.Context, child dream.ChildhoodDream) error {
	input := transform.NewInput(child)

	output, err := r.useCase.Execute(ctx, input)
	if err != nil {
		return appErrors.Wrap(err, appErrors.CodeDomainFailure, "useCase execution failed")
	}

	vm, err := r.presenter.Present(ctx, output)
	if err != nil {
		return appErrors.Wrap(err, appErrors.CodeInternal, "presenter failed")
	}

	text, err := r.formatter.Format(ctx, vm)
	if err != nil {
		return appErrors.Wrap(err, appErrors.CodeInternal, "formatter failed")
	}

	if _, err := fmt.Fprint(r.out, text); err != nil {
		return appErrors.Wrap(err, appErrors.CodeIO, "failed to write formatted output")
	}

	return nil
}
