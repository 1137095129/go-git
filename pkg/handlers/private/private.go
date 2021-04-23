package private

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/sirupsen/logrus"
	"github.com/wang1137095129/go-git/config"
	"github.com/wang1137095129/go-git/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Private struct {
	localpath string
}

func (p *Private) OpenRepository(c *config.Config) (*git.Repository, error) {
	if len(p.localpath) > 0 {
		return git.PlainOpen(p.localpath)
	} else {
		p.localpath = filepath.Join(utils.HomeDir(), c.Git.RepositoryName)
		file, err := os.Open(c.User.CertificatePath)
		if err != nil {
			logrus.Fatal(err)
		}
		defer file.Close()
		b, err := ioutil.ReadAll(file)
		if err != nil {
			logrus.Fatal(err)
		}
		publicKeys, err := ssh.NewPublicKeys(c.User.Username, b, c.User.Password)
		if err != nil {
			logrus.Fatal(err)
		}
		return git.PlainClone(
			p.localpath,
			false,
			&git.CloneOptions{
				URL:  c.Git.URL,
				Auth: publicKeys,
				ReferenceName: plumbing.ReferenceName(c.Git.Branch) ,
			},
		)
	}
}

func (p *Private) Refresh(c *config.Config) (*git.Repository, error) {
	repository, err := p.OpenRepository(c)
	if err != nil {
		logrus.Fatal(err)
	}
	worktree, err := repository.Worktree()
	if err != nil {
		logrus.Fatal(err)
	}
	file, err := os.Open(c.User.CertificatePath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Fatal(err)
	}
	publicKeys, err := ssh.NewPublicKeys(c.User.Username, b, c.User.Password)
	if err != nil {
		logrus.Fatal(err)
	}
	err = worktree.Pull(&git.PullOptions{
		Auth: publicKeys,
	})
	return repository, nil
}
