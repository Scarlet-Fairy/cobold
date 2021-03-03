package git

import (
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	gitUtils "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

type gitClone struct {
}

func newClone() clone.Clone {
	return &gitClone{}
}

func (g *gitClone) Clone(_ context.Context, options clone.Options) error {
	_, err := gitUtils.PlainClone(options.Path, false, &gitUtils.CloneOptions{
		URL:               options.Url,
		RecurseSubmodules: gitUtils.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return errors.Wrap(err, "could not clone repo")
	}

	return nil
}
