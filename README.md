# Hetzanetes

Create K3s Kubernetes clusters on Hetzner Cloud. With apologies to Hetzner and Kubernetes for the name! 

## What does it do?

Right now, nothing. There is no code - just a plan. The aim is to:

* Provide a simple way to set up and manage Kubernetes clusters on Hetzner Cloud.
* Avoid local configuration files.
* Work with Rancher's lightweight [K3s](https://github.com/rancher/k3s/) Kubernetes distribution. 
* Install Hetzner's [cloud controller manager](https://github.com/hetznercloud/hcloud-cloud-controller-manager) and [storage volume](https://github.com/hetznercloud/csi-driver) plugins.
* Set up a firewall and private network for the cluster, like [Vito Botta](https://github.com/vitobotta/hetzner-cloud-init) does.
* Automate security updates.
* Make the cluster as self-managing as possible.

## Alternatives

I wanted a simple way to create and manage Kubernetes clusters on Hetzner Cloud. There are really good projects out there but none of them quite did what I wanted (as of 2020-09-19).
They are all worth checking out, especially if this project doesn't meet your needs.

* [Pharmer](https://github.com/pharmer/pharmer) - loads of features but doesn't support Hetzner Cloud. Kubernetes versions are a bit old.
* [Hetzner-Kube](https://github.com/xetys/hetzner-kube) - impressive networking setup dates from before Hetzner Cloud had private networks and load balancers. Uses `kubeadm`.
* [K3sup](https://github.com/alexellis/k3sup) - great way to install Rancher's K3s Kubernetes on a cluster but it doesn't actually provision the cluster and doesn't set up the firewall K3s needs.

## How it works

1. Run `hetzanetes` locally. If this is the first run, you need to provide a [Hetzner Cloud API token](https://console.hetzner.cloud/projects) > (your project) > Security > API tokens.
2. Hetzanetes will set up a private network in your Hetzner Cloud project, attach some labels for configuration purposes, and create a single Kubernetes API server node.
3. When the single API server node starts, it will create further API server and worker nodes and join them to the cluster.
4. All nodes are joined to the private network, have a firewall (using `ufw`) and are set up for unattended upgrades.
