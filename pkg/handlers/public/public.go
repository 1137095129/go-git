package public

import (
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"github.com/wang1137095129/go-git/config"
	"github.com/wang1137095129/go-git/utils"
	"os"
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
		_, err := os.Stat(p.localpath)
		if err != nil {
			if os.IsNotExist(err) {
				return git.PlainClone(p.localpath, false, &git.CloneOptions{
					URL:           c.Git.URL,
					RemoteName:    c.Git.Branch,
				})
			}
			return nil, err
		}
		return git.PlainOpen(p.localpath)
	}
}

func (p *Public) Refresh(c *config.Config) (*git.Repository, error) {
	repository, err := git.PlainOpenWithOptions(
		p.localpath,
		&git.PlainOpenOptions{
			EnableDotGitCommonDir: true,
			DetectDotGit:          true,
		},
	)
	if err != nil {
		logrus.Fatal(err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		logrus.Fatal(err)
	}
	err = worktree.Pull(&git.PullOptions{
		RemoteName: c.Git.Branch,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return repository, nil
}
