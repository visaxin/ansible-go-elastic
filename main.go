package main

import (
	"flag"
	"os"

	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-ansible-elastic-cluster/bootstrap"
	"github.com/go-ansible-elastic-cluster/handler"
)

func main() {

	once := sync.Once{}
	once.Do(bootstrap.TemplateFile)
	var listenPort string

	flag.StringVar(&listenPort, "p", os.Getenv("PORT_HTTP"), "listen port")
	flag.Parse()

	if listenPort == "" {
		listenPort = "8080"
	}
	router := gin.Default()

	//v1
	v1 := router.Group("/api/v1")
	{
		v1.POST("/cluster", handler.CoreHandler)
		v1.POST("/deploy", handler.DeployHandler)
		v1.GET("/deploy/status", handler.DeployStatusHandler)
		v1.GET("/deploys", handler.DeployListHandler)

		// plugin
		v1.GET("/plugins", handler.PluginHandler)
		v1.POST("/plugins", handler.UploadPlugin)
	}
	os.Setenv("S_PORT", listenPort)
	router.Run(":" + listenPort)

	return
}
