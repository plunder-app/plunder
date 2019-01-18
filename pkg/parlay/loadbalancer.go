package parlay

import (
	"fmt"
	"os/exec"
)

type loadBalancer struct {
	// Load balancer
	LBHostname        string `json:"lbHostname,omitempty"`
	LBPort            int    `json:"lbPort,omitempty"`
	EndPointHostnames []string
	EndPointAddresses []string
	EndPointPorts     []int

	// Defines the type of load balancer and if it will be ran in a container
	LBRType     string `json:"lbType"`
	LBContainer bool   `json:"lbContainer,omitempty"`
	LBImage     string `json:"lbImage,omitempty"`
}

//HA Proxy configuration file outline
const haproxyCfg = `frontend k8s-api
    bind %s:%d
    bind 127.0.0.1:6443
    mode tcp
    option tcplog
    default_backend k8s-api

backend k8s-api
    mode tcp
    option tcplog
    option tcp-check
    balance roundrobin
    default-server inter 10s downinter 5s rise 2 fall 2 slowstart 60s maxconn 250 maxqueue 256 weight 100

`

func (l *loadBalancer) buildHAProxycfg() (string, error) {
	var cfg string
	if len(l.EndPointAddresses) != len(l.EndPointHostnames) {
		return "", fmt.Errorf("Endpoint address count doesn't match hostnames")
	}
	if len(l.EndPointAddresses) != len(l.EndPointPorts) {
		return "", fmt.Errorf("Endpoint address count doesn't match ports")
	}
	cfg = fmt.Sprintf(haproxyCfg, l.LBHostname, l.LBPort)

	for x := range l.EndPointHostnames {
		cfg = fmt.Sprintf("%s    server %s %s:%d check\n", cfg, l.EndPointHostnames[x], l.EndPointAddresses[x], l.EndPointPorts[x])
	}

	// Generates a haproxy configuration file
	return cfg, nil
}

// This function will generate all of the required actions to configure a load balancer
func (l *loadBalancer) generateActions() ([]Action, error) {
	var generatedActions []Action
	var a Action
	var configString, configPath string
	var err error

	// This will generate the required configuration files needed for a specific load balancer
	switch l.LBRType {
	case "haproxy":
		configString, err = l.buildHAProxycfg()
		configPath = "/etc/haproxy/haproxy.cfg"
		if err != nil {
			return nil, err
		}
	case "nginx":
		// TODO
	default:
		return nil, fmt.Errorf("Unknown Load balancer type [%s]", l.LBRType)
	}

	a = Action{
		// Configure the load balancer
		ActionType:  "command",
		Name:        fmt.Sprintf("Configure the load balancer [%s]", l.LBRType),
		Command:     fmt.Sprintf("echo '%s' > %s", configString, configPath),
		CommandSudo: "root",
	}

	generatedActions = append(generatedActions, a)

	if l.LBContainer == true {
		// Shell out, TODO - perhaps hit up the docker socket
		exec.Command("`which docker`")
	}

	return generatedActions, nil
}
