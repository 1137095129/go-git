package main

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"github.com/libgit2/git2go/v31"
	"os/exec"
	"time"
)

func main() {
	timer := time.NewTimer(1 * time.Minute)
	for true {
		<-timer.C
		timer.Reset(1 * time.Minute)
		now := time.Now()
		ago := now.Add(-5 * time.Minute)
		go func() {
			repository, err := git.OpenRepository("https://github.com/1137095129/springboot-netty.git")
			if err!=nil {
				logrus.Fatal(err)
				return
			}
			oid := &git.Oid{}
			commit, _ := repository.LookupCommit(oid)
			when := commit.Committer().When
			fmt.Println(fmt.Sprintf("%d-%d-%d %d:%d:%d", when.Year(), when.Month(), when.Day(), when.Hour(), when.Minute(), when.Second()))
			fmt.Println(fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))
			fmt.Println(fmt.Sprintf("%d-%d-%d %d:%d:%d", ago.Year(), ago.Month(), ago.Day(), ago.Hour(), ago.Minute(), ago.Second()))
			output, _ := exec.Command("git", "log", "--help").Output()
			bytess, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(output)
			buffer := bytes.NewBuffer(bytess)
			fmt.Println(buffer.String())
		}()
	}
	//logrus.Entry{}
}
