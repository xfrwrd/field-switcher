package presenter

import (
	"context"
	"fmt"
	"slices"

	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	"github.com/xeniasokk/field-switcher/internal/ports"
	appErrors "github.com/xeniasokk/field-switcher/pkg/errors"
)

var _ ports.PresenterPort = (*ConsolePresenter)(nil)
var _ ports.ViewModel = ConsoleViewModel{}

type ConsoleViewModel struct {
	title     string
	childhood dream.ChildhoodDream
	adult     dream.Adult
	note      string
}

func (vm ConsoleViewModel) Title() string {
	return vm.title
}

func (vm ConsoleViewModel) Childhood() dream.ChildhoodDream {
	return vm.childhood
}

func (vm ConsoleViewModel) Adult() dream.Adult {
	return vm.adult
}

func (vm ConsoleViewModel) Note() string {
	return vm.note
}

func NewConsoleViewModel(
	title string,
	childhood dream.ChildhoodDream,
	adult dream.Adult,
	note string,
) ConsoleViewModel {
	return ConsoleViewModel{
		title:     title,
		childhood: childhood,
		adult:     adult,
		note:      note,
	}
}

type ConsolePresenter struct{}

func NewConsolePresenter() *ConsolePresenter {
	return &ConsolePresenter{}
}

func (p *ConsolePresenter) Present(
	ctx context.Context,
	output ports.OutputModel,
) (ports.ViewModel, error) {
	_ = ctx

	adult := output.Adult()
	if adult.RoleTitle() == "" {
		return ConsoleViewModel{}, appErrors.NewDomainError("adult identity has empty role title")
	}

	child := output.Child()
	if child.DisplayName() == "" {
		return ConsoleViewModel{}, appErrors.NewDomainError("childhood dream has empty display name")
	}

	// Проверяем наличие упорства
	traits := adult.Traits()
	hasPersistence := slices.ContainsFunc(traits, func(q dream.Quality) bool {
		return q.Name() == dream.QualityPersistence
	})

	note := fmt.Sprintf("Сохранено качеств: %d", len(adult.Traits()))
	if hasPersistence {
		note += fmt.Sprintf(" | %s — твой главный союзник на новом поле", dream.QualityPersistence)
	}

	vm := NewConsoleViewModel(
		"field-switcher — трансформация мечты",
		child,
		adult,
		note,
	)

	return vm, nil
}
