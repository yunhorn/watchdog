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

	client = dingtalk.NewClient(*token, *secret)

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		fmt.Println("hello world")
		wf := workflow{}
		err := c.BindJSON(&wf)
		if err != nil {
			log.Println("parse.request.failed!", err)
		} else {

			if wf.Action != "" {
				log.Println("action is", wf.Action)
				if wf.workflow_run != nil {
					content := fmt.Sprintf("Github Action Job:%s,headBranch:%s,status:%s.\r\n", wf.workflow_run.name, wf.workflow_run.head_branch, wf.workflow_run.status)
					msg := dingtalk.NewTextMessage().SetContent(content).SetAt([]string{"mobile", ""}, false)
					client.Send(msg)
					log.Println("action wf name:", wf.workflow_run.name)
					log.Println("action wf status:", wf.workflow_run.status)
					log.Println("action wf head_branch:", wf.workflow_run.head_branch)
				}
			}
		}

		c.String(200, "success")
	})
	router.Run(":8181")
}
