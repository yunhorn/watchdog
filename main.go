package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/gin-gonic/gin"
)

type workflow struct {
	Action      string      `json:"action"`
	WorkflowRun workflowRun `json:"workflow_run,omitempty"`
}

type workflowRun struct {
	Name         string     `json:"name,omitempty"`
	HeadBranch   string     `json:"head_branch"`
	Status       string     `json:"status"`
	Conclusion   string     `json:"conclusion"`
	HtmlUrl      string     `json:"html_url"`
	RunStartedAt string     `json:"run_started_at"`
	Head_commit  headCommit `json:"head_commit"`
}

type headCommit struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
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
		var wf workflow
		err := c.BindJSON(&wf)
		if err != nil {
			log.Println("parse.request.failed!", err)
		} else {

			if wf.Action != "" {
				log.Println("action is", wf.Action)
				log.Println("workflow run is nil", wf.WorkflowRun)
				if wf.WorkflowRun.Name != "" {
					content := fmt.Sprintf("Github Action Job:%s,headBranch:%s,status:%s.\r\n", wf.WorkflowRun.Name, wf.WorkflowRun.HeadBranch, wf.WorkflowRun.Status)
					msg := dingtalk.NewTextMessage().SetContent(content).SetAt([]string{"mobile", ""}, false)
					client.Send(msg)
				}
			}
		}

		c.String(200, "success:"+wf.WorkflowRun.Name)
	})
	router.Run(":8181")
}
