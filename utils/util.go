package utils

import (
	"os"
	"runtime"
)

func HomeDir() string {
	if home := os.Getenv("GO_GIT_HOME"); home != "" {
		return home
	}
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}