package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ansible-elastic-cluster/core"
	"github.com/go-ansible-elastic-cluster/metadata"
)

func CoreHandler(c *gin.Context) {
	request := c.Request
	cluster := &core.Cluster{}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	err = json.Unmarshal(body, cluster)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err})
	}

	// TODO register resource request
	metadata.DataSource(nil).Save()
	// TODO start to deploy

	// TODO confirm deploy result

	// TODO release resource when deploy fail

	c.JSON(200, gin.H{})

}
