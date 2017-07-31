package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ansible-elastic-cluster/core"
)

func DeployHandler(c *gin.Context) {
	name := c.Query("name")
	go func() {
		stdout, err := core.ExecuteDeploy(name)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error(), "msg": string(stdout)})
			return
		}
	}()
	c.JSON(200, gin.H{"msg": string("submit success!")})
}

func DeployStatusHandler(c *gin.Context) {
	name := c.Query("name")

	status, err := core.DeployStatus(name)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
			"msg":   string(status),
		})
		return
	}
	c.String(200, "%s", string(status))
}

func DeployListHandler(c *gin.Context) {
	name := c.Query("name")

	list, err := core.DeployList(name)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
			"msg":   "fail to list cluster name " + name,
		})
		return
	}
	c.JSON(200, gin.H{
		"deploy_history": list,
	})
}
