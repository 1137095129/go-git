package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wang1137095129/go-git/config"
	"github.com/wang1137095129/go-git/pkg/client"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "go_git",
	Short: "watch git repository.",
	Long: `
watch git repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		c := &config.Config{}
		if err := c.Load(); err != nil {
			logrus.Fatal(err)
		}
		fmt.Println("starting...")
		client.Run(c)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		c := &config.Config{}
		if err := c.Load(); err != nil {
			logrus.Fatal(err)
		}
		fmt.Println("starting...")
		client.Run(c)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigName(config.ConfigFileName)
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file", viper.ConfigFileUsed())
	}
}
