package core

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
)

type Instance struct {
	Config      map[string]interface{}
	DataPathDir string // "dir1,dir2"
	LogPathDir  string
}

type Host struct {
	Instances []Instance
	HostName  string
}

type Plugin struct {
	Name    string
	Version string
	Url     string
}

func (this *Plugin) SetUrl() {
	this.Url = fmt.Sprintf("http://example.com/%s_%s.zip", this.Name, this.Version)
}

type Cluster struct {
	Hosts       []Host
	ClusterName string
	EsConfig    map[string]interface{} // for some common config in a cluster
	JVMConfig   map[string]interface{} // for config jvm
	Mode        string
	DataPathDir []string
	LogPathDir  string
	Plugins     []Plugin
}

func Create(name string, hosts []Host, dataPathDir []string, logPathDir string) *Cluster {
	c := Cluster{
		Hosts:       hosts,
		ClusterName: name,
		DataPathDir: dataPathDir,
		LogPathDir:  logPathDir,
		JVMConfig:   map[string]interface{}{"es_heap_size": "20g"},
	}

	initInstanceConfig(c)
	zenPingList := zenPingList(c)
	esConfig := make(map[string]interface{}, 0)
	esConfig["discovery.zen.ping.unicast.hosts"] = strings.Join(zenPingList, ",")
	for k, v := range commonConfig {
		esConfig[k] = v
	}
	c.EsConfig = esConfig

	return &c
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

// init for some special config
// node.maser
// http.port
// transport.tcp.port
func initInstanceConfig(c Cluster) {
	hosts := c.Hosts
	// standalone mode: just use runtime.
	// multi-node cluster:  use metadata from agent
	var cpu int = 2
	if c.Mode == StandaloneMode {
		cpu = runtime.NumCPU()
	}
	for sh, h := range hosts {
		processorMax := cpu / len(h.Instances)
		if processorMax < 1 {
			processorMax = 1
		}
		for serial, i := range h.Instances {
			i.Config = make(map[string]interface{}, 0)

			var httpPort = 9200 + serial
			var transPort = 9300 + serial
			if _, ok := i.Config["http.port"]; !ok {
				i.Config["http.port"] = httpPort
			}
			if _, ok := i.Config["transport.tcp.port"]; !ok {
				i.Config["transport.tcp.port"] = transPort
			}

			i.Config["processors"] = processorMax
			i.Config["node.master"] = true
			i.Config["network.host"] = h.HostName
			i.Config["network.publish_host"] = h.HostName
			i.DataPathDir = strings.Join(dataPathAllocation(c.DataPathDir, serial, len(h.Instances)), ",")
			i.LogPathDir = c.LogPathDir
			h.Instances[serial] = i
		}
		hosts[sh] = h
	}
}

func dataPathAllocation(dataPathDir []string, index int, total int) []string {
	avg := len(dataPathDir) / total
	start := index * avg
	return dataPathDir[start : start+avg]
}

func (this *Cluster) generateAnsibleYml() error {
	var err error
	var parsedTemplate *template.Template
	parsedTemplate, err = template.ParseFiles("deploy.yml")
	if err != nil {
		return err
	}

	templateBuff := new(bytes.Buffer)
	err = parsedTemplate.Execute(templateBuff, this)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s-s/deploy.yml", this.ClusterName, time.Now().Format(time.RFC3339))
	return ioutil.WriteFile(fileName, templateBuff.Bytes(), 0655)
}

func (this *Cluster) Run() {
	this.generateAnsibleYml()
}
