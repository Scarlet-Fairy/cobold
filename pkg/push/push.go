package push

import (
	"context"
)

const (
	StepName = "push"
)

type Push interface {
	Push(ctx context.Context, options Options) error
}

type Options struct {
	Name     string
	Tag      string
	Registry string
}
