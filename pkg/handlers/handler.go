package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/wang1137095129/go-git/config"
)

type Handler interface {
	OpenRepository(c *config.Config) (*git.Repository, error)
	Refresh(c *config.Config) (*git.Repository, error)
}
