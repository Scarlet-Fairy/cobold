package build

import (
	"context"
	"io"
)

const (
	StepName = "build"
)

type Build interface {
	Build(ctx context.Context, options Options) (io.Reader, error)
}

type Options struct {
	Directory string
	Name      string
}
