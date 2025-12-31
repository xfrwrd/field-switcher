package ports

import "context"

type PresenterPort interface {
	Present(ctx context.Context, output OutputModel) (ViewModel, error)
}
