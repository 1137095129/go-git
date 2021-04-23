package client

import (
	"github.com/wang1137095129/go-git/config"
	"github.com/wang1137095129/go-git/pkg/controller"
	"github.com/wang1137095129/go-git/pkg/handlers"
	"github.com/wang1137095129/go-git/pkg/handlers/private"
	"github.com/wang1137095129/go-git/pkg/handlers/public"
)

func Run(c *config.Config) {
	controller.Start(c, getHandlerByConfig(c))
}

func getHandlerByConfig(c *config.Config) handlers.Handler {
	if len(c.User.CertificatePath) > 0 && len(c.User.Password) > 0 && len(c.User.Username) > 0 {
		return new(private.Private)
	}
	return new(public.Public)
}
