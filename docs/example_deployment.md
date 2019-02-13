# Example Deployment for off-line Kubernetes

**This example will make use of Plunders User Interface**

In order for an offline installation to be succesful, a lot of the packages and containers will need downloading to where Plunder will be ran from. 

## Offline Calico parts

### Download the manifests

```
wget https://docs.projectcalico.org/v3.5/getting-started/kubernetes/installation/hosted/etcd.yaml
```
and

```
wget https://docs.projectcalico.org/v3.5/getting-started/kubernetes/installation/hosted/calico.yaml
```

### Download the named images

One Liner to pull the calico images and etcd image

``` 
for image in $(cat etcd.yaml | grep image | awk '{ print $2 }') ; do sudo docker pull $image; done
```
``` 
for image in $(cat calico.yaml | grep image | awk '{ print $2 }') ; do sudo docker pull $image; done
```

At this point you'll have the images as part of the local docker repository and the two manifests in the local directory.

## Offline Ubuntu packages

One liner to get teh packages needed for the kubernetes hosts to run `kubelet`

```
apt-get download socat ethtool ebtables; tar -cvzf ubuntu_pkg.tar.gz socat* ethtool* ebtables*; rm socat* ethtool* ebtables*
```

This command will download everything needed into an archive named `ubuntu_pkg.tar.gz`

## Offline Docker packages 

One liner to get the docker-ce packages for all kubernetes hosts, ensure that the docker repository has been added to the hosts repositories before attempting to download the package. 

```
apt-get download docker-ce=18.06.1~ce~3-0~ubuntu; tar -cvzf docker_pkg.tar.gz ./docker-ce_18.06.1~ce~3-0~ubuntu_amd64.deb; rm docker-ce_18.06.1~ce~3-0~ubuntu_amd64.deb
```

## Offline Kubernetes packages

One liner to get the kubernetes packages for all kubernetes hosts, ensure that the kubernets repository has been added to the hosts repositories before attempting to download the package. 

```
apt-get download kubelet kubeadm kubectl cri-tools kubernetes-cni; tar -cvzf kubernetes_pkg.tar.gz kubelet* kubeadm* kubectl* cri-tools* kubernetes-cni*; rm kubelet* kubeadm* kubectl* cri-tools* kubernetes-cni*
```

## Offline Kubernetes images

The easiest way of managing this is to install kubeadm on the pluder host and use `kubeadm` to prep the local docker image store with the images needed.

`kubeadm config images list` - will list all images

`kubeadm config images pull` - will pull them all to the local host

Once all of the images have been pulled locally or downloaded as tars manually from the registry we can modify out deployment map and deploy as expected. 

## Example deployment map

There is an example deployment map as a `gist` available [https://gist.github.com/thebsdbox/f12b621a9d3943128b6bb16688497cd0](https://gist.github.com/thebsdbox/f12b621a9d3943128b6bb16688497cd0)

## Deployment in action

[![asciicast](https://asciinema.org/a/reh3reEgJQKCOB5e92D96l6tt.png)](https://asciinema.org/a/reh3reEgJQKCOB5e92D96l6tt)

