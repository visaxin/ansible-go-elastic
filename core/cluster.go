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
	if this.EsConfig == nil {
		this.EsConfig = make(map[string]interface{})
	}
	initInstanceConfig(*this)
	this.EsConfig["discovery.zen.ping.unicast.hosts"] = strings.Join(zenPingList(*this), ",")
	for k, v := range commonConfig {
		this.EsConfig[k] = v
	}
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

	path = fmt.Sprintf("%s/%s-%d", cacheDir, this.ClusterName, time.Now().Unix())
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return path, err
	}

	fileName := fmt.Sprintf("%s/%s", path, DefaultYmlFile)
	return path, ioutil.WriteFile(fileName, templateBuff.Bytes(), 0755)
}

func (this *Cluster) CreateConfigFile() (string, error) {
	return this.generateAnsibleYml(DefaultCacheDir, DefaultYmlFile)

}
