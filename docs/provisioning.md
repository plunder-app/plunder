# Provisioning Configuration

The provisioning works by running remote commands or uploading/downloading files to a remote system, in order for it to be configured correctly. A parsing engine called "parlay" was written in order to provide repeatable scripting to ease deployments.

A Deployment map can contain multiple **deployments**, which in turn will contain one or more **actions** that will be performed on one or more **hosts**.

Also a deployment map can be parsed as either **JSON** or as **YAML** (yaml being somewhat easier to read as a human and creating much smaller files).

### Example deployment map

This script below (for offline installations) will upload a tarball containing the docker packages and then install them on all remote systems listed under `hosts`.

**Note** the tarball was created by `apt-get download docker-ce=18.06.1~ce~3-0~ubuntu; tar -cvzf docker_pkg.tar.gz ./docker-ce_18.06.1~ce~3-0~ubuntu_amd64.deb`

#### JSON Example

```json
{
	"deployments": [
		{
			"name": "Upload Docker Packages",
			"parallel": false,
			"sessions": 0,
			"hosts": [
				"192.168.1.3",
				"192.168.1.4",
				"192.168.1.5"
			],
			"actions": [
				{
					"name": "Upload Docker Packages",
					"type": "upload",
					"source": "./docker_pkg.tar.gz",
					"destination": "/tmp/docker_pkg.tar.gz"
				},
				{
					"name": "Extract Docker packages",
					"type": "command",
					"command": "tar -C /tmp -xvzf /tmp/docker_pkg.tar.gz"
				},
        {
					"name": "Install Docker packages",
					"type": "command",
					"command": "dpkg -i /tmp/docker/*",
					"commandSudo": "root"
				}
			]
		}
	]
}
                                
```
#### YAML Example

```yaml
deployments:
- actions:
  - destination: /tmp/docker_pkg.tar.gz
    name: Upload Docker Packages
    source: ./docker_pkg.tar.gz
    timeout: 0
    type: upload
  - command: tar -C /tmp -xvzf /tmp/docker_pkg.tar.gz
    name: Extract Docker packages
    timeout: 0
    type: command
  - command: dpkg -i /tmp/docker/*
    commandSudo: root
    name: Install Docker packages
    timeout: 0
    type: command
  hosts:
  - 192.168.1.3
  - 192.168.1.4
  - 192.168.1.5
  name: Upload Docker Packages
  parallel: false
  parallelSessions: 0
```

The above example only covers simple usage of `uploading` and `command` usages.



## Usage

When automating a deployment ssh credentials are required to map a host with the correct credentials. To simplify this `plunder` can make use of the deployment file to determine access credentials. When the deployment begins `plunder` will evaluate the hosts in the provisioning map and identify the correct credentials from the deployment file. 

`plunder automate ssh --config ./deployment.json --map ./etcd.json --logfile output`

![](../image/parlay.jpg)
*The above example uses screen, where the output from `plunder` is on the top and `tail -f output` is below*

Additional flags:

- The `--deployment` flag now will point to a specific deployment in a map
- The `--action` flag can be used to point to a specific action in a deployment
- The `--host` flag will point to a specific host in the deployment
- The `--resume` will determine if to continue executing all remaining actions
