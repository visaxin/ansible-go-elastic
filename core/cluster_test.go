package core

import (
	"encoding/json"
	"io/ioutil"
	"testing"

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

	exceptESScript := false
	assert.Equal(t, exceptESScript, c.Vars["es_scripts"])

	exceptList := []string{"host19:9300", "host19:9301", "host19:9302"}
	list := zenPingList(*c)
	assert.Equal(t, 9, len(list))
	assert.Equal(t, exceptList, list[:3])
	exceptLogPath := "/disk1/log"
	assert.Equal(t, exceptLogPath, c.Hosts[0].Instances[0].LogPathDir)

	defer os.RemoveAll(DefaultCacheDir)
	_, err = c.generateAnsibleYml(DefaultCacheDir, "../deploy.yml")
	assert.NoError(t, err)

	var name string
	fs, err := ioutil.ReadDir(DefaultCacheDir)
	for _, f := range fs {
		name = f.Name()
	}

	_, err = ExecuteDeploy(name)
	assert.Error(t, err)

	// test get execute status
	status, err := DeployStatus(name)
	assert.NoError(t, err)
	t.Log(string(status))

	cl, err := listDeployHistory(exceptClusterName)

	assert.NoError(t, err)
	assert.Equal(t, name, cl[0])

}
