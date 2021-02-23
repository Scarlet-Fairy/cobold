package build

import "io"

type Build interface {
	Build(options BuildOptions) (io.Reader, error)
}

type BuildOptions struct {
	Directory string
	Tag       string
}
