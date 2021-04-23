package config

import (
	"github.com/wang1137095129/go-git/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ConfigFileName = ".gogit.yaml"
)

type Config struct {
	Git  Git  `json:"git" yaml:"git"`
	User User `json:"user" yaml:"user"`
}

type Git struct {
	Branch         string `json:"branch" yaml:"branch,omitempty"`
	RepositoryName string `json:"repositoryName" yaml:"repositoryName,omitempty"`
	URL            string `json:"url" yaml:"url,omitempty"`
}

type User struct {
	Username        string `json:"username" yaml:"username,omitempty"`
	Password        string `json:"password" yaml:"password,omitempty"`
	CertificatePath string `json:"certificatePath" yaml:"certificatePath,omitempty"`
}

func New() (*Config, error) {
	c := &Config{}
	if err := c.Load(); err != nil {
		return c, err
	}
	return c, nil
}

func (c *Config) Load() error {
	err := createIfNotExists()
	if err != nil {
		return err
	}
	file, err := os.Open(getConfigFile())
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) > 0 {
		return yaml.Unmarshal(b, c)
	}
	return nil
}

func (c *Config) Write() error {
	file, err := os.OpenFile(getConfigFile(), os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	enc.SetIndent(2)
	return enc.Encode(c)
}

func createIfNotExists() error {
	filePath := filepath.Join(utils.HomeDir(), ConfigFileName)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filePath)
			if err != nil {
				return err
			}
			file.Close()
		} else {
			return err
		}
	}
	return nil
}

func getConfigFile() string {
	filePath := filepath.Join(utils.HomeDir(), ConfigFileName)
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	}
	return ""
}
