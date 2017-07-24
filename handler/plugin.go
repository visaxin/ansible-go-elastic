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
