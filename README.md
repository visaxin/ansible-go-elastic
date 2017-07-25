# ansible-go-elastic
One-click to build production level Elasticsearch


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
          "host_name": "cs19"
        },
        {
          "instances": [
            {},
            {},
            {}
          ],
          "host_name": "cs20"
        },
        {
          "instances": [
            {},
            {},
            {}
          ],
          "host_name": "cs21"
        }
      ],
      "cluster_name": "test",
      "data_path_dir": [
        "/disk1",
        "/disk2",
        "/disk3",
        "/disk4",
        "/disk5"
      ],
      "log_path_dir": "/disk1/log"
    }
    '
    
    
    * {} represents an Instance. You can provide a new config：

    	{
    		"log_path_dir": "/var/log/",
    		"data_path_dir": "/var/data/",
    		"config":{
    			"bootstrap.mlockall": false
    		}
    	}