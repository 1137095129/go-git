package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/wang1137095129/go-git/config"
)

var gitConfigCmd = &cobra.Command{
	Use:   "gitconfig",
	Short: "specific git-repository configuration",
	Long:  `specific git-repository configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}
		url, err := cmd.Flags().GetString("url")
		if err == nil {
			if len(url) > 0 {
				c.Git.URL = url
			}
		} else {
			logrus.Fatal(err)
		}
		branch, err := cmd.Flags().GetString("branch")
		if err == nil {
			if len(branch) > 0 {
				c.Git.Branch = branch
			}
		} else {
			logrus.Fatal(err)
		}
		username, err := cmd.Flags().GetString("username")
		if err == nil {
			if len(username) > 0 {
				c.User.Username = username
			}
		} else {
			logrus.Fatal(err)
		}
		password, err := cmd.Flags().GetString("password")
		if err == nil {
			if len(password) > 0 {
				c.User.Password = password
			}
		} else {
			logrus.Fatal(err)
		}
		certificate, err := cmd.Flags().GetString("certificate")
		if err == nil {
			if len(certificate) > 0 {
				c.User.CertificatePath = password
			}
		} else {
			logrus.Fatal(err)
		}
		repository, err := cmd.Flags().GetString("repository")
		if err == nil {
			if len(repository) > 0 {
				c.Git.RepositoryName = repository
			}
		}else {
			logrus.Fatal(err)
		}
		remote, err := cmd.Flags().GetString("remote")
		if err == nil {
			if len(remote) > 0 {
				c.Git.RemoteName = remote
			}
		}else {
			logrus.Fatal(err)
		}
		if err = c.Write(); err != nil {
			logrus.Fatal(err)
		}
	},
}

func init() {
	gitConfigCmd.Flags().StringP("url", "", "", "Specify git-repository url")
	gitConfigCmd.Flags().StringP("branch", "b", "master", "Specify git-branch")
	gitConfigCmd.Flags().StringP("username", "u", "", "Specify git-username")
	gitConfigCmd.Flags().StringP("password", "p", "", "Specify git-password")
	gitConfigCmd.Flags().StringP("certificate", "c", "", "Specify git-certificate")
	gitConfigCmd.Flags().StringP("repository", "r", "default", "Specify git-repository-name")
	gitConfigCmd.Flags().StringP("remote","","origin","Specify git-remote-name")
}
