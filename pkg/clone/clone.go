package clone

import "context"

type Clone interface {
	Clone(ctx context.Context, url string, path string) error
}
