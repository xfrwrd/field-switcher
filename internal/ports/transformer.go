package ports

import (
	"context"

	"github.com/xeniasokk/field-switcher/internal/domain/dream"
)

type TransformerPort interface {
	TransformDream(ctx context.Context, d dream.ChildhoodDream) (dream.Adult, error)
}
