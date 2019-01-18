package parlay

import (
	"fmt"
)

// NOTE - The functions in this particluar file will need moving to something seperate at a later
// date. Quite possibly moving to a plugin model? TBD.

// This defines the etcd kubeadm file (should use the kubernetes packages to define at a later point)
const etcdKubeadm = `apiVersion: "kubeadm.k8s.io/%s"
kind: ClusterConfiguration
etcd:
    local:
        serverCertSANs:
        - "%s"
        peerCertSANs:
        - "%s"
        extraArgs:
            initial-cluster: %s=https://%s:2380,%s=https://%s:2380,%s=https://%s:2380
            initial-cluster-state: new
            name: %s
            listen-peer-urls: https://%s:2380
            listen-client-urls: https://%s:2379
            advertise-client-urls: https://%s:2379
            initial-advertise-peer-urls: https://%s:2380`

// This defines the manager kubeadm file (should use the kubernetes packages to define at a later point)

const managerKubeadm = `apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
kubernetesVersion: %s
apiServer:
  certSANs:
  - "%s"
controlPlaneEndpoint: "%s:%d"
etcd:
    external:
        endpoints:
        - https://%s:2379
        - https://%s:2379
        - https://%s:2379
        caFile: /etc/kubernetes/pki/etcd/ca.crt
        certFile: /etc/kubernetes/pki/apiserver-etcd-client.crt
        keyFile: /etc/kubernetes/pki/apiserver-etcd-client.key`

type etcdMembers struct {
	// Hostnames
	Hostname1 string `json:"hostname1,omitempty"`
	Hostname2 string `json:"hostname2,omitempty"`
	Hostname3 string `json:"hostname3,omitempty"`

	// Addresses
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	Address3 string `json:"address3,omitempty"`

	// Intialise a Certificate Authority
	InitCA bool `json:"initCA,omitempty"`

	// Set kubernetes API version
	APIVersion string `json:"apiversion,omitempty"`
}

type managerMembers struct {
	// ETCD Nodes
	ETCDAddress1 string `json:"etcd01,omitempty"`
	ETCDAddress2 string `json:"etcd02,omitempty"`
	ETCDAddress3 string `json:"etcd03,omitempty"`

	// Version of Kubernetes
	Version string `json:"kubeVersion,omitempty"`

	// Load Balancer details (needed for initialising the first master)
	loadBalancer

	// Stacked - means ETCD nodes are stacked on managers (false by default)
	Stacked bool `json:"stacked,omitempty"`
}

func (e *etcdMembers) generateActions() []Action {
	var generatedActions []Action
	var a Action
	if e.InitCA == true {
		// Ensure that a new Certificate Authority is generated
		// Create action
		a = Action{
			// Generate etcd server certificate
			ActionType:  "command",
			Command:     fmt.Sprintf("kubeadm init phase certs etcd-ca"),
			CommandSudo: "root",
			Name:        "Initialise Certificate Authority",
		}
		generatedActions = append(generatedActions, a)
	}

	// Default to < 1.12 API version
	if e.APIVersion == "" {
		e.APIVersion = "v1beta1"
	}
	// Generate the configuration directories
	a.ActionType = "command"
	a.Command = fmt.Sprintf("mkdir -m 777 -p /tmp/%s/ /tmp/%s/ /tmp/%s/", e.Address1, e.Address2, e.Address3)
	a.Name = "Generate temporary directories"
	generatedActions = append(generatedActions, a)

	// Generate the kubeadm configuration files

	// Node 0
	a.Name = "build kubeadm config for node 0"
	a.Command = fmt.Sprintf("echo '%s' > /tmp/%s/kubeadmcfg.yaml", e.buildKubeadm(e.APIVersion, e.Hostname1, e.Address1), e.Address1)
	generatedActions = append(generatedActions, a)

	// Node 1
	a.Name = "build kubeadm config for node 1"
	a.Command = fmt.Sprintf("echo '%s' > /tmp/%s/kubeadmcfg.yaml", e.buildKubeadm(e.APIVersion, e.Hostname2, e.Address2), e.Address2)
	generatedActions = append(generatedActions, a)

	// Node 2
	a.Name = "build kubeadm config for node 2"
	a.Command = fmt.Sprintf("echo '%s' > /tmp/%s/kubeadmcfg.yaml", e.buildKubeadm(e.APIVersion, e.Hostname3, e.Address3), e.Address3)
	generatedActions = append(generatedActions, a)

	// Add certificate actions
	generatedActions = append(generatedActions, e.generateCertificateActions([]string{e.Address3, e.Address2, e.Address1})...)
	return generatedActions
}

func (e *etcdMembers) buildKubeadm(api, host, address string) string {
	var kubeadm string
	// Generates a kubeadm for setting up the etcd yaml
	kubeadm = fmt.Sprintf(etcdKubeadm, api, address, address, e.Hostname1, e.Address1, e.Hostname2, e.Address2, e.Hostname3, e.Address3, host, address, address, address, address)
	return kubeadm
}

// generateCertificateActions - Hosts need adding in backward to the array i.e. host 2 -> host 1 -> host 0
func (e *etcdMembers) generateCertificateActions(hosts []string) []Action {
	var generatedActions []Action
	var a Action

	a.Command = "mkdir -p /etc/kubernetes/pki"
	a.CommandSudo = "root"
	a.Name = "Ensure that PKI directory exists"
	a.ActionType = "command"
	generatedActions = append(generatedActions, a)

	for i, v := range hosts {
		// Tidy any existing client certificates
		a.ActionType = "command"
		a.Command = "find /etc/kubernetes/pki -not -name ca.crt -not -name ca.key -type f -delete"
		a.Name = "Remove any existing client certificates before attempting to generate any new ones"
		generatedActions = append(generatedActions, a)

		// Generate etcd server certificate
		a.ActionType = "command"
		a.Command = fmt.Sprintf("kubeadm init phase certs etcd-server --config=/tmp/%s/kubeadmcfg.yaml", v)
		a.Name = fmt.Sprintf("Generate etcd server certificate for [%s]", v)
		generatedActions = append(generatedActions, a)

		// Generate peer certificate
		a.Command = fmt.Sprintf("kubeadm init phase certs etcd-peer --config=/tmp/%s/kubeadmcfg.yaml", v)
		a.Name = fmt.Sprintf("Generate peer certificate for [%s]", v)
		generatedActions = append(generatedActions, a)

		// Generate health check certificate
		a.Command = fmt.Sprintf("kubeadm init phase certs etcd-healthcheck-client --config=/tmp/%s/kubeadmcfg.yaml", v)
		a.Name = fmt.Sprintf("Generate health check certificate for [%s]", v)
		generatedActions = append(generatedActions, a)

		// Generate api-server client certificate
		a.Command = fmt.Sprintf("kubeadm init phase certs apiserver-etcd-client --config=/tmp/%s/kubeadmcfg.yaml", v)
		a.Name = fmt.Sprintf("Generate api-server client certificate for [%s]", v)
		generatedActions = append(generatedActions, a)

		// These steps are only required for the first two hosts
		if i != (len(hosts) - 1) {
			// Archive the certificates and the kubeadm configuration in a host specific archive name
			a.Command = fmt.Sprintf("tar -cvzf /tmp/%s.tar.gz $(find /etc/kubernetes/pki -type f) /tmp/%s/kubeadmcfg.yaml", v, v)
			a.Name = fmt.Sprintf("Archive generated certificates [%s]", v)
			generatedActions = append(generatedActions, a)

			// Download the archive files to the local machine
			a.ActionType = "download"
			a.Source = fmt.Sprintf("/tmp/%s.tar.gz", hosts[i])
			a.Destination = fmt.Sprintf("/tmp/%s.tar.gz", hosts[i])
			a.Name = fmt.Sprintf("Retrieve the certificate bundle for [%s]", v)
			generatedActions = append(generatedActions, a)
		} else {
			// This is the final host, grab the certificates for use by a manager
			a.Command = fmt.Sprintf("tar -cvzf /tmp/managercert.tar.gz /etc/kubernetes/pki/etcd/ca.crt /etc/kubernetes/pki/apiserver-etcd-client.crt /etc/kubernetes/pki/apiserver-etcd-client.key")
			a.Name = fmt.Sprintf("Archive generated certificates [%s]", v)
			generatedActions = append(generatedActions, a)

			// Download the archive files to the local machine
			a.ActionType = "download"
			a.Source = "/tmp/managercert.tar.gz"
			a.Destination = "/tmp/managercert.tar.gz"
			a.Name = "Retrieving the Certificates for the manager nodes"
			generatedActions = append(generatedActions, a)
		}
	}
	return generatedActions
}

// At some point the functions for the various kubeadm arease will be split into seperate files to ease management
func (m *managerMembers) generateActions() []Action {
	var generatedActions []Action
	var a Action
	if m.Stacked == false {
		// Not implemented yet TODO
		return nil
	}

	// Upload the initial etcd certificates to the first manager node
	a = Action{
		// Upload etcd server certificate
		ActionType:  "upload",
		Source:      "/tmp/managercert.tar.gz",
		Destination: "/tmp/managercert.tar.gz",
		Name:        "Upload etcd server certificate to first manager",
	}
	generatedActions = append(generatedActions, a)

	// Install the certificates for etcd
	a.Name = "Installing the etcd certificates"
	a.ActionType = "command"
	a.CommandSudo = "root"
	a.Command = fmt.Sprintf("tar -xvzf /tmp/managercert.tar.gz -C /")
	generatedActions = append(generatedActions, a)

	// Generate the kubeadm configuration file
	a.Name = "Generating the Kubeadm file for the first manager node"
	a.Command = fmt.Sprintf("echo '%s' > /tmp/kubeadmcfg.yaml", m.buildKubeadm())
	generatedActions = append(generatedActions, a)

	// Initialise the first node
	a.Name = "Initialise the first control plane node"
	a.Command = "kubeadm init --config /tmp/kubeadmcfg.yaml"
	generatedActions = append(generatedActions, a)

	return generatedActions
}

func (m *managerMembers) buildKubeadm() string {
	var kubeadm string
	// Generates a kubeadm for setting up the etcd yaml
	kubeadm = fmt.Sprintf(managerKubeadm, m.Version, m.LBHostname, m.LBHostname, m.LBPort, m.ETCDAddress1, m.ETCDAddress2, m.ETCDAddress3)
	return kubeadm
}
