package ports

import "github.com/xeniasokk/field-switcher/internal/domain/dream"

type ViewModel interface {
	Title() string
	Childhood() dream.ChildhoodDream
	Adult() dream.Adult
	Note() string
}
