package core

// cluster level config
var commonConfig = map[string]interface{}{
	"gateway.recover_after_nodes":        2,
	"discovery.zen.ping_timeout":         "10s",
	"discovery.zen.minimum_master_nodes": 1,
	"bootstrap.mlockall":                 true,
	"indices.fielddata.cache.size":       "30%",
	"indices.breaker.fielddata.limit":    "40%",
	"format": "json",
	"index.translog.flush_threshold_size":   "3g",
	"index.translog.flush_threshold_period": "40m",
	"index.translog.interval":               "10s",
	"indices.store.throttle.type":           "none",
	"action.write_consistency":              "one",
	"security.manager.enabled":              "false",
	"script.groovy.sandbox.enabled":         "true",
	"script.inline":                         "on",
	"script.indexed":                        "on",
	"node.max_local_storage_nodes":          1,
	"node.data":                             true,
}

const (
	StandaloneMode string = "standalone"

	DefaultCacheDir = ".go-ansible"
	DefaultYmlFile  = "deploy.yml"
)
