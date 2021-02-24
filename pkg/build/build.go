package build

import (
	"context"
	"io"
)

type Build interface {
	Build(ctx context.Context, options Options) (io.Reader, error)
}

type Options struct {
	Directory string
	Tag       string
}
