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

		//X-GitHub-Event: workflow_run/workflow_job
		// not we just work for workflow_run
		value := c.GetHeader("X-GitHub-Event")
		if value != "workflow_run" {
			c.String(200, "success")
			return
		}

		var wf workflow
		err := c.BindJSON(&wf)
		if err != nil {
			log.Println("parse.request.failed!", err)
		} else {

			if wf.Action != "" {
				log.Println("action is", wf.Action)
				log.Println("workflow run is nil", wf.WorkflowRun)
				if wf.WorkflowRun.Name != "" {
					content := fmt.Sprintf("Github Action:%s,\n\n 分支:%s, \n\n状态:%s.\r \n\n最新提交信息:%s \n\n[详情](%s)", wf.WorkflowRun.Name, wf.WorkflowRun.HeadBranch, wf.WorkflowRun.Status, wf.WorkflowRun.Head_commit.Message, wf.WorkflowRun.HtmlUrl)
					msg := dingtalk.NewMarkdownMessage().SetMarkdown("action通知", content).SetAt([]string{"mobile", ""}, false)
					client.Send(msg)
				}
			}
		}

		c.String(200, "success:"+wf.WorkflowRun.Name)
	})
	router.Run(":8181")
}
