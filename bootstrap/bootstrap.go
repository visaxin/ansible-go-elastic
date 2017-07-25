package bootstrap

import (
	"io/ioutil"

	"log"

	"github.com/go-ansible-elastic-cluster/core"
)

var template = []byte(`---
{{ range $h := $.Hosts }}{{ range $i,$e := $h.Instances }}
-
  hosts: {{ $h.HostName }}
  roles:
    -
      es_config:
        cluster.name: {{$.ClusterName}}
        {{ range $k,$v := $e.Config }}{{ $k }}: "{{ $v }}"
        {{ end }}
      es_heap_size: {{index $.JVMConfig "es_heap_size" }}
      es_instance_name: {{ $h.HostName }}-{{ $i }}
      es_data_dirs: "{{ $e.DataPathDir }}"
      es_log_dir: "{{ $e.LogPathDir }}"
      role: elasticsearch
  vars:
    es_plugins:
      -
        {{ range $p := $.Plugins }}plugin: {{ $p.Url }} {{ end }}
    {{ range $k,$v := $.Vars }}{{ $k }}: "{{ $v }}"
    {{ end }}
{{ end }}{{ end }}
`)

// generate a template yml file using go-template syntax
func TemplateFile() {
	err := ioutil.WriteFile(core.DefaultYmlFile, template, 0655)
	if err != nil {
		log.Fatal(err)
	}
}
