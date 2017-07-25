package core

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"
)

type Instance struct {
	Config      map[string]interface{} `json:"config"`
	DataPathDir string                 `json:"data_path_dir" ` // "dir1,dir2"
	LogPathDir  string                 `json:"log_path_dir"`
}

type Host struct {
	Instances []Instance `json:"instances"`
	HostName  string     `json:"host_name"`
}

type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	url     string
}

func (this *Plugin) SetUrl() string {
	hostName, _ := os.Hostname()
	// TODO how to connect to localhost:<PORT>
	this.url = fmt.Sprintf("http://%s/%s-%s.zip", hostName, this.Name, this.Version)
	return this.url
}

type Cluster struct {
	Hosts       []Host                 `json:"hosts"`
	ClusterName string                 `json:"cluster_name"`
	EsConfig    map[string]interface{} `json:"es_config"`  // for some common config in a cluster
	JVMConfig   map[string]interface{} `json:"jvm_config"` // for config jvm
	Mode        string                 `json:"mode"`
	DataPathDir []string               `json:"data_path_dir"`
	LogPathDir  string                 `json:"log_path_dir"`
	Plugins     []Plugin               `json:"plugins"`
}

func (this *Cluster) Init() *Cluster {
	initInstanceConfig(this)
	updateInstanceConfig(this)
	return this
}

func Create(name string, hosts []Host, dataPathDir []string, logPathDir string) *Cluster {
	c := Cluster{
		Hosts:       hosts,
		ClusterName: name,
		DataPathDir: dataPathDir,
		LogPathDir:  logPathDir,
		JVMConfig:   map[string]interface{}{"es_heap_size": "20g"},
	}
	return c.Init()
}

// select master node from Hosts
// require: initInstanceConfig(hosts)
func zenPingList(c Cluster) []string {
	hosts := c.Hosts
	zenPingList := []string{}
	for _, h := range hosts {
		for _, i := range h.Instances {
			isMaster, ok := i.Config["node.master"]
			if ok && isMaster == true {
				connUrl := fmt.Sprintf("%s:%d", h.HostName, i.Config["transport.tcp.port"])
				zenPingList = append(zenPingList, connUrl)
			}
		}
	}
	return zenPingList
}

func updateInstanceConfig(c *Cluster) {
	list := strings.Join(zenPingList(*c), ",")
	for sh, h := range c.Hosts {
		for serial, i := range h.Instances {
			i.Config["discovery.zen.ping.unicast.hosts"] = list
			h.Instances[serial] = i
		}
		c.Hosts[sh] = h
	}
}

// init for some special config
// node.maser
// http.port
// transport.tcp.port
func initInstanceConfig(c *Cluster) {
	// standalone mode: just use runtime.
	// multi-node cluster:  use metadata from agent
	var cpu int = 2
	if c.Mode == StandaloneMode {
		cpu = runtime.NumCPU()
	}
	c.JVMConfig = map[string]interface{}{"es_heap_size": "5g"}
	for sh, h := range c.Hosts {
		processorMax := cpu / len(h.Instances)
		if processorMax < 1 {
			processorMax = 1
		}
		for serial, i := range h.Instances {
			config := make(map[string]interface{})
			for k, v := range commonConfig {
				config[k] = v
			}
			for k, v := range i.Config {
				config[k] = v
			}
			var httpPort = 9200 + serial
			var transPort = 9300 + serial
			config["http.port"] = httpPort
			config["transport.tcp.port"] = transPort

			if _, ok := config["processors"]; !ok {
				config["processors"] = processorMax
			}
			if _, ok := config["node.master"]; !ok {
				config["node.master"] = true
			}

			config["network.host"] = h.HostName
			config["network.publish_host"] = h.HostName

			if i.DataPathDir == "" {
				i.DataPathDir = strings.Join(dataPathAllocation(c.DataPathDir, serial, len(h.Instances)), ",")
			}
			if i.LogPathDir == "" {
				i.LogPathDir = c.LogPathDir
			}
			i.Config = config
			h.Instances[serial] = i
			c.Hosts[sh] = h
		}
	}
}

func dataPathAllocation(dataPathDir []string, index int, total int) []string {
	avg := len(dataPathDir) / total
	start := index * avg
	return dataPathDir[start : start+avg]
}

func (this *Cluster) generateAnsibleYml(cacheDir string, templateFile string) (string, error) {

	var err error
	var path string
	var parsedTemplate *template.Template
	parsedTemplate, err = template.ParseFiles(templateFile)
	if err != nil {
		return path, err
	}

	templateBuff := new(bytes.Buffer)
	err = parsedTemplate.Execute(templateBuff, this)
	if err != nil {
		return path, err
	}

	uid := fmt.Sprintf("%s-%d", this.ClusterName, time.Now().Unix())
	path = fmt.Sprintf("%s/%s", cacheDir, uid)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return path, err
	}

	fileName := fmt.Sprintf("%s/%s", path, DefaultYmlFile)
	return uid, ioutil.WriteFile(fileName, templateBuff.Bytes(), 0755)
}

func (this *Cluster) CreateConfigFile() (string, error) {
	return this.generateAnsibleYml(DefaultCacheDir, DefaultYmlFile)

}
