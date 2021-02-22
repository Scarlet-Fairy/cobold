package git

import (
	"context"
	"github.com/Scarlet-Fairy/cobold/pkg/clone"
	gitUtils "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
)

type git struct {
}

func new() clone.Clone {
	return &git{}
}

func (g *git) Clone(_ context.Context, url string, path string) error {
	_, err := gitUtils.PlainClone(path, false, &gitUtils.CloneOptions{
		URL:               url,
		RecurseSubmodules: gitUtils.DefaultSubmoduleRecursionDepth,
	})
	if err != nil {
		return errors.Wrap(err, "git.Clone")
	}

	return nil
}
