package transform

import (
	"github.com/xeniasokk/field-switcher/internal/domain/dream"
	"github.com/xeniasokk/field-switcher/internal/ports"
)

var _ ports.OutputModel = OutputModel{}

type OutputModel struct {
	child dream.ChildhoodDream
	adult dream.Adult
}

func NewOutput(child dream.ChildhoodDream, adult dream.Adult) OutputModel {
	return OutputModel{
		child: child,
		adult: adult,
	}
}

func (o OutputModel) Child() dream.ChildhoodDream { return o.child }
func (o OutputModel) Adult() dream.Adult          { return o.adult }
