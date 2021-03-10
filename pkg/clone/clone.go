package clone

import "context"

const (
	StepName = "clone"
)

type Clone interface {
	Clone(ctx context.Context, options Options) error
}

type Options struct {
	Url  string
	Path string
}
