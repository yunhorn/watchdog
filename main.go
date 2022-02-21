package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type workflow struct {
	action       string
	workflow_run workflowRun
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

func main() {
	router := gin.Default()
	router.GET("/webhook", func(c *gin.Context) {
		fmt.Println("hello world")
		c.String(200, "success")
	})
	router.Run(":8181")
}
