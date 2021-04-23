package controller

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"github.com/wang1137095129/go-git/config"
	"github.com/wang1137095129/go-git/pkg/handlers"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var once = &sync.Once{}

type Controller struct {
	lastPullTime *time.Time
}

func Start(conf *config.Config, eventHandler handlers.Handler) {
	timer := time.NewTimer(1 * time.Minute)
	c := timer.C

	controller := &Controller{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGINT)
	for true {
		select {
		case <-c:
			timer.Reset(1 * time.Minute)
			go controller.Run(conf, eventHandler)
		case <-signals:
			fmt.Println("stop controller")
			return
		}
	}
}

func (c *Controller) Run(conf *config.Config, eventHandler handlers.Handler) {
	var flag = false
	once.Do(func() {
		c.lastPullTime = new(time.Time)
		_, err := eventHandler.OpenRepository(conf)
		if err != nil {
			logrus.Fatal(err)
		}
		*c.lastPullTime = time.Now()
		flag = true
	})

	if flag {
		return
	}

	repository, err := eventHandler.OpenRepository(conf)
	if err != nil {
		logrus.Fatal(err)
	}
	head, err := repository.Head()
	if err != nil {
		logrus.Fatal(err)
	}

	commitIter, err := repository.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		From:  head.Hash(),
		Since: c.lastPullTime,
	})

	if err != nil {
		logrus.Fatal(err)
	}

	var t time.Time

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Committer.When.After(t) {
			t = commit.Committer.When
		}
		return nil
	})

	if c.lastPullTime.Before(t) {
		*c.lastPullTime = t
		_, err := eventHandler.Refresh(conf)
		if err!=nil {
			logrus.Fatal(err)
		}
	}

}
