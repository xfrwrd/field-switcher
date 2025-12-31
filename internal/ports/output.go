package ports

import "github.com/xeniasokk/field-switcher/internal/domain/dream"

type OutputModel interface {
	Child() dream.ChildhoodDream
	Adult() dream.Adult
}
