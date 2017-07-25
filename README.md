# ansible-go-elastic
One-click to build production level Elasticsearch.

Dependency: https://github.com/elastic/ansible-elasticsearch


Example:

    Config a cluster:
    
    curl -XPOST localhost:8080/api/v1/cluster --data '
    {
      "hosts": [
        {
          "instances": [
            {},
            {},
            {}
          ],
          "host_name": "host19"
        },
        {
          "instances": [
            {},
            {},
            {}
          ],
          "host_name": "host20"
        },
        {
          "instances": [
            {},
            {},
            {}
          ],
          "host_name": "host21"
        }
      ],
      "cluster_name": "test",
      "data_path_dir": [
        "/path1",
        "/path2",
        "/path3",
        "/path4",
        "/path5"
      ],
      "log_path_dir": "/path/to/log"
    }
    '
    
    
 * {} represents an Instance. You can provide a new config for one instance：

    	{
    		"log_path_dir": "/var/log/",
    		"data_path_dir": "/var/data/",
    		"config":{
    			"bootstrap.mlockall": false
    		}
    	}
