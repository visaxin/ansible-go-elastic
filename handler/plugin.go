package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func PluginHandler(c *gin.Context) {
	name := c.Query("name")
	version := c.Query("version")

	fileName := fmt.Sprintf("/tmp/%s-%s.zip", name, version)
	c.File(fileName)
}

func UploadPlugin(c *gin.Context) {
	name := c.Query("name")
	version := c.Query("version")

	fileName := fmt.Sprintf("/tmp/%s-%s.zip", name, version)
	file, err := c.FormFile("plugin.zip")
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "cannot get upload file"})
		return
	}
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "cannot save upload file"})
		return
	}
}
