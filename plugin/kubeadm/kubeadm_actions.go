package main

import (
	"fmt"

	"github.com/plunder-app/plunder/pkg/parlay/parlaytypes"
)

func (e *etcdMembers) generateActions() []parlaytypes.Action {
	var generatedActions []parlaytypes.Action
	var a parlaytypes.Action
	if e.InitCA == true {
		// Ensure that a new Certificate Authority is generated
		// Create action
		a = parlaytypes.Action{
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
func (e *etcdMembers) generateCertificateActions(hosts []string) []parlaytypes.Action {
	var generatedActions []parlaytypes.Action
	var a parlaytypes.Action

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
func (m *managerMembers) generateActions() []parlaytypes.Action {
	var generatedActions []parlaytypes.Action
	var a parlaytypes.Action
	if m.Stacked == false {
		// Not implemented yet TODO
		return nil
	}

	// Upload the initial etcd certificates to the first manager node
	a = parlaytypes.Action{
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
	kubeadm = fmt.Sprintf(managerKubeadm, m.Version, "LB HOSTNAME FIXME", "LB HOSTNAME FIXME", 1000000, m.ETCDAddress1, m.ETCDAddress2, m.ETCDAddress3)
	return kubeadm
}
