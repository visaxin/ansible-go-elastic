# ansible-go-elastic
One-click to build production level Elasticsearch.

## Setup guide
> require
*  Go 1.8.x 
*  Ansible 2.x


1. Install ansible 

2. Clone ansible-elasticsearch and make `elasticsearch` role availiable in ansible globlly.
  
      `git clone https://github.com/elastic/ansible-elasticsearch`
      `sudo mkdir -p /etc/ansible/roles/ && mv ansible-elasticsearch /etc/ansible/roles`
  
3.    `git clone https://github.com/visaxin/ansible-go-elastic.git && cd ansible-go-elastic && go build`

4. start `go-ansible-elastic-cluster` default listen on `8080`




## Example

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
