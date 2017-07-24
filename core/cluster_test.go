package core

import "testing"
import (
	"encoding/json"
	"io/ioutil"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestDataPathAllocation(t *testing.T) {
	path := []string{"/path1", "/path2", "/path3"}

	allocation := dataPathAllocation(path, 2, 3)
	assert.Equal(t, "/path3", allocation[0])

	path = []string{"/path1", "/path2", "/path3", "/path4", "/path5", "/path6", "/path7", "/path8"}

	allocation = dataPathAllocation(path, 1, 2)
	except := []string{"/path5", "/path6", "/path7", "/path8"}
	assert.Equal(t, except, allocation)
}

func TestClusterInput(t *testing.T) {
	b, err := ioutil.ReadFile("example.json")
	assert.NoError(t, err)
	c := &Cluster{}
	err = json.Unmarshal(b, c)
	assert.NoError(t, err)

	c = c.Init()

	exceptClusterName := "test"
	assert.Equal(t, exceptClusterName, c.ClusterName)

	exceptLogPath := "/disk1/log"
	assert.Equal(t, exceptLogPath, c.Hosts[0].Instances[0].LogPathDir)

	defer os.RemoveAll(".cache")
	assert.NoError(t, c.generateAnsibleYml(".cache", "../deploy.yml"))
}
