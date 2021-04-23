package controller

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

var lock = &sync.WaitGroup{}

type Controller struct {
	lastPullTime *time.Time
}

func Start(conf *config.Config, eventHandler handlers.Handler) {
	timer := time.NewTimer(0 * time.Minute)
	c := timer.C

	controller := &Controller{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGINT)
	for true {
		select {
		case <-c:
			timer.Reset(10 * time.Second)
			go controller.Run(conf, eventHandler)
		case <-signals:
			fmt.Println("stop controller")
			return
		}
	}
}

func (c *Controller) Run(conf *config.Config, eventHandler handlers.Handler) {
	lock.Wait()
	var flag = false
	fmt.Println(fmt.Sprintf("check git url: %s ,branch name:%s,repository name:%s", conf.Git.URL, conf.Git.Branch, conf.Git.RepositoryName))
	//刚开始初始化
	once.Do(func() {
		c.lastPullTime = new(time.Time)
		fmt.Println("init repository")
		r, err := eventHandler.OpenRepository(conf)
		if err != nil {
			logrus.Fatal(err)
		}
		head, err := r.Head()
		if err != nil {
			logrus.Fatal(err)
		}
		commitIter, err := r.Log(&git.LogOptions{From: head.Hash()})
		if err != nil {
			logrus.Fatal(err)
		}
		defer commitIter.Close()
		var t time.Time
		commitIter.ForEach(func(commit *object.Commit) error {
			if t.Before(commit.Committer.When) {
				t = commit.Committer.When
			}
			return nil
		})
		if c.lastPullTime.Before(t) {
			*c.lastPullTime = t
			fmt.Println(
				fmt.Sprintf(
					"get last commit when %d-%d-%d %d:%d:%d",
					t.Year(),
					t.Month(),
					t.Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
				),
			)
		}
		flag = true
	})

	if flag {
		return
	}

	lock.Wait()
	lock.Add(1)
	defer lock.Done()

	//打开git仓库
	repository, err := eventHandler.OpenRepository(conf)

	if err != nil {
		logrus.Fatal(err)
	}

	//从远端获取更新
	repository.Fetch(&git.FetchOptions{RemoteName: conf.Git.Branch})

	fmt.Println(fmt.Sprintf("check %s/%s", conf.Git.RemoteName, conf.Git.Branch))

	//获取远端最新一次更新
	resolveRevision, err := repository.ResolveRevision(plumbing.Revision(fmt.Sprintf("%s/%s", conf.Git.RemoteName, conf.Git.Branch)))
	if err != nil {
		logrus.Fatal(err)
	}
	commit, err := repository.CommitObject(*resolveRevision)
	if err != nil {
		logrus.Fatal(err)
	}

	//更新git仓库
	if c.lastPullTime.Before(commit.Committer.When) {
		*c.lastPullTime = commit.Committer.When
		fmt.Println(
			fmt.Sprintf(
				"get new commit when %d-%d-%d %d:%d:%d, will be refresh",
				c.lastPullTime.Year(),
				c.lastPullTime.Month(),
				c.lastPullTime.Day(),
				c.lastPullTime.Hour(),
				c.lastPullTime.Minute(),
				c.lastPullTime.Second(),
			),
		)
		_, err := eventHandler.Refresh(conf)
		if err != nil {
			logrus.Fatal(err)
		}
		fmt.Println("pull success")
	}

}
