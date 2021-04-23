package public

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"github.com/wang1137095129/go-git/config"
	"github.com/wang1137095129/go-git/utils"
	"path/filepath"
)

type Public struct {
	localpath string
}

func (p *Public) OpenRepository(c *config.Config) (*git.Repository, error) {
	if len(p.localpath) > 0 {
		return git.PlainOpen(p.localpath)
	} else {
		p.localpath = filepath.Join(utils.HomeDir(), c.Git.RepositoryName)
		return git.PlainClone(p.localpath, false, &git.CloneOptions{
			URL:           c.Git.URL,
			ReferenceName: plumbing.ReferenceName(c.Git.Branch),
		})
	}
}

func (p *Public) Refresh(c *config.Config) (*git.Repository, error) {
	repository, err := p.OpenRepository(c)
	if err != nil {
		logrus.Fatal(err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		logrus.Fatal(err)
	}
	err = worktree.Pull(&git.PullOptions{

	})
	if err != nil {
		logrus.Fatal(err)
	}
	return repository, nil
}
