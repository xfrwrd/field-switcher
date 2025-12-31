package app

import (
	"context"
	"os"

	"github.com/xeniasokk/field-switcher/internal/adapters/formatter"
	"github.com/xeniasokk/field-switcher/internal/adapters/presenter"
	"github.com/xeniasokk/field-switcher/internal/adapters/runner"
	"github.com/xeniasokk/field-switcher/internal/adapters/transformer"
	"github.com/xeniasokk/field-switcher/internal/application/usecase/transform"
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	appErrors "github.com/xeniasokk/field-switcher/pkg/errors"
	"github.com/xeniasokk/field-switcher/pkg/lifecycle"
)

type app struct {
	runner *runner.ConsoleRunner
	dream  dream.ChildhoodDream
}

func NewApp() (lifecycle.App, error) {
	return NewAppWithConfig(DefaultConfig())
}

func NewAppWithConfig(cfg Config) (lifecycle.App, error) {
	tr, err := transformer.NewSimpleTransformer(cfg.TargetRole)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.CodeInternal, "create transformer")
	}

	uc := transform.NewUseCase(tr)

	p := presenter.NewConsolePresenter()
	f := formatter.NewTextFormatter()

	r, err := runner.NewConsoleRunner(uc, p, f, os.Stdout)
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.CodeInternal, "create console runner")
	}

	defaultDream, err := dream.NewDefaultFootballerDream()
	if err != nil {
		return nil, appErrors.Wrap(err, appErrors.CodeInternal, "create default footballer dream")
	}

	return &app{
		runner: r,
		dream:  defaultDream,
	}, nil
}

func (a *app) Run(ctx context.Context) error {
	return a.runner.Run(ctx, a.dream)
}

func (a *app) Shutdown(ctx context.Context) error {
	_ = ctx
	return nil
}
