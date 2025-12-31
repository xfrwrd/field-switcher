package ports

import "context"

type FormatterPort interface {
	Format(ctx context.Context, vm ViewModel) (string, error)
}
