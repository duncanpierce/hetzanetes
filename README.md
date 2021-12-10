# Hetzanetes

Create K3s Kubernetes clusters on Hetzner Cloud. With apologies to Hetzner and Kubernetes for the name! 

## What does it do?

Right now, it can only provision a single server with K3s Kubernetes. The aim is to:

* Provide a simple way to set up and manage Kubernetes clusters on Hetzner Cloud.
* Avoid local configuration files; be able to manage clusters from anywhere, provided you have an API token.
* Work with Rancher's lightweight [K3s Kubernetes distribution](https://github.com/rancher/k3s/). 
* Avoid proliferation of installation options, OS base images, etc.
* Install Hetzner's [cloud controller manager](https://github.com/hetznercloud/hcloud-cloud-controller-manager) and [storage volume](https://github.com/hetznercloud/csi-driver) plugins.
* Set up a firewall and private network for the cluster, like [Vito Botta](https://github.com/vitobotta/hetzner-cloud-init) does.
* Automate security updates, where possible.
* Make the cluster as self-repairing as possible.
* Synchronize SSH keys the cluster will accept with those registered in the Hetzner API - handy if your lose you private key or your team changes.
* Make SSH recognise new Hetzner servers so we don't get "key changed" errors.

## Alternatives

I wanted a simple way to create and manage Kubernetes clusters on Hetzner Cloud. There are really good projects out there but none of them quite did what I wanted (as of 2020-09-19).
They are all worth checking out, especially if this project doesn't meet your needs.

* [Pharmer](https://github.com/pharmer/pharmer) - loads of features but doesn't support Hetzner Cloud.
* [Hetzner-Kube](https://github.com/xetys/hetzner-kube) - impressive networking setup dates from before Hetzner Cloud had private networks, load balancers and labels. Uses `kubeadm`.
* [K3sup](https://github.com/alexellis/k3sup) - great way to install Rancher's K3s Kubernetes on a cluster but it doesn't provision the cluster or up a firewall.
* [Kube-Hetzner](https://github.com/mysticaltech/kube-hetzner) - uses Terraform to set up K3OS
* [kubernetes-on-hetzner](https://github.com/LWJ/kubernetes-on-hetzner) - uses Terraform

## How it works

1. `hetzanetes` used to share configuration with [Hetzner's command line interpreter](https://github.com/hetznercloud/cli) but it's now decoupled, so you need to set an environment variable HCLOUD_TOKEN with an API token obtained from [Hetzner Cloud API token](https://console.hetzner.cloud/projects) > (your project) > Security > API tokens.
2. Hetzanetes will set up a private network in your Hetzner Cloud project, attach some labels for configuration purposes, and create a single Kubernetes API server node.
3. (in future) When the single API server node starts, it will create further API server and worker nodes and join them to the cluster.
4. All nodes are joined to the private network, have a firewall (using `ufw`) and are set up for unattended upgrades.
