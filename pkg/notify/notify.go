package notify

import "context"

type Notify interface {
	NotifyCompletion(ctx context.Context, options Options) error
}

type Options struct {
	JobID  string
	Reason string
	Err    error
}
