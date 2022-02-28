package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/gin-gonic/gin"
)

type workflow struct {
	Action       string
	workflow_run *workflowRun
}

type workflowRun struct {
	name           string
	head_branch    string
	status         string
	conclusion     string
	html_url       string
	run_started_at string
	head_commit    headCommit
}

type headCommit struct {
	message   string
	timestamp string
}

var (
	client *dingtalk.Client
)

func main() {

	token := flag.String("token", "", "token for dingtalk")
	secret := flag.String("secret", "", "secret for dingtalk")
	flag.Parse()

	if *token == "" || *secret == "" {
		log.Fatalln("please set token or secret")
	}

	// https://oapi.dingtalk.com/robot/send?access_token=
	// accessToken := "3df22b6d3f3d3a9fdb4e8df846d995b594c7610aee4699069d4900dcad834996"
	// secret := "SEC4bc4ab892a8e23f16e6de0fd92b9ac5a0fe92a8fe63d06a74e79b3233c17c0d4"
	client = dingtalk.NewClient(*token, *secret)

	router := gin.Default()
	router.GET("/webhook", func(c *gin.Context) {
		fmt.Println("hello world")
		wf := workflow{}
		err := c.BindJSON(&wf)
		if err != nil {
			log.Println("parse.request.failed!", err)
		} else {

			log.Println("action is", wf.Action)
			if wf.workflow_run != nil {
				log.Println("action wf name:", wf.workflow_run.name)
				log.Println("action wf status:", wf.workflow_run.status)
				log.Println("action wf head_branch:", wf.workflow_run.head_branch)
			}
		}
		// msg := dingtalk.NewTextMessage().SetContent("测试文本&at 某个人").SetAt([]string{"mobile", ""}, false)
		// client.Send(msg)
		c.String(200, "success")
	})
	router.Run(":8181")
}
