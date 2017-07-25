package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ansible-elastic-cluster/core"
)

func DeployHandler(c *gin.Context) {
	name := c.Query("name")
	stdout, err := core.ExecuteDeploy(name)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"msg": string(stdout)})
}
