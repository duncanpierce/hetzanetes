[![CodeQL](https://github.com/duncanpierce/hetzanetes/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/duncanpierce/hetzanetes/actions/workflows/github-code-scanning/codeql)

# Hetzanetes

Create self-managing [Rancher](https://www.rancher.com/products/k3s) [K3s](https://github.com/k3s-io/k3s) Kubernetes clusters on [Hetzner Cloud](https://www.hetzner.com/cloud).

* The cluster manages itself using a `Cluster` [Custom Resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
* Avoids local configuration files. You can reconfigure the cluster using `kubectl edit cluster/<NAME>`. The cluster can even reconfigure itself.
* Self-repairing, so if a node is deleted, Hetzanetes will replace it.
* Number of API servers and workers can be changed without stopping the cluster.
* Sets up Rancher's lightweight [K3s Kubernetes distribution](https://github.com/rancher/k3s/).
* Sets up a firewall and private network for the cluster.
* Install Hetzner's [cloud controller manager](https://github.com/hetznercloud/hcloud-cloud-controller-manager) and [storage volume](https://github.com/hetznercloud/csi-driver) plugins, so volume and load balancer resources work.

## Current limitations

* Downsizing the API server node set can hang the cluster. Worker node pools can be downsized.
* Cannot manage other clusters, even though you could have more than one `Cluster` resource, in theory.

## Getting started

1. Create a Hetzner Cloud project, if you don't already have one. If this is your first time using Hetzner Cloud, you can get [free service credit using this referral link](https://hetzner.cloud/?ref=0obWH96gyXNL) (disclosure: I may also receive a service credit).
2. Create a read+write API Token in that project (under **Security > API Tokens**), if you don't already have one.
3. Assign the API Token to an environment variable named `HCLOUD_TOKEN`.
4. Run `hetzanetes create test` to create a cluster called `test`. Or, using Docker: `docker run -e HCLOUD_TOKEN ghcr.io/duncanpierce/hetzanetes create test`.
5. Wait patiently while a private network, firewall and first Kubernetes API server are created, security updates are installed, server rebooted, Hetzner's Kubernetes plugins installed.
6. Once ready, the first API server will read the `Cluster` custom resource and create more API servers and workers as needed.
7. From this point on, the cluster is self-managing. The complete process takes around 10 minutes using CX11 servers.
8. You can now log into any of the API servers and use `kubectl edit cluster/test` (or whatever cluster name you chose) to reconfigure the cluster.

## Configuring the cluster before creation

Save this example as a file named `cluster.yaml` then customise it with the server nodes you want.
This example has a single set of API server nodes and two sets of worker nodes.

```yaml
apiVersion: hetzanetes.duncanpierce.org/v1
kind: Cluster
metadata:
  name: example
spec:
  versions:
    baseImage: ubuntu-22.04
    kubernetes: stable
    hetzanetes: latest
  nodeSets:
    - name: api
      apiServer: true
      serverType: cpx21
      replicas: 3
      locations:
        - hel1
        - fsn1
        - nbg1
    - name: worker-pool-a
      serverType: cpx11
      replicas: 2
      locations:
        - hel1
        - fsn1
    - name: worker-pool-b
      serverType: cpx11
      replicas: 1
      locations:
        - nbg1
```

Then run `hetzanetes create -f cluster.yaml` to create the first API server node. This will then create the rest of the cluster.

## In future

* Release prebuild executables to avoid building hetzanetes yourself.
* Automate security updates.
* Automate K3s distribution updates.
* Synchronize SSH keys the cluster will accept with those registered in the Hetzner API - handy if your lose you private key or your team changes.
* Make SSH recognise new Hetzner servers so we don't get "key changed" errors.
* Include workloads at creation time to be run in the cluster once it's ready.
* Make it easy to download the kube config file.
* Optionally create a load balancer for the API servers to make it easier to use `kubectl` remotely.

## Alternatives

I wanted a simple way to create and manage Kubernetes clusters on Hetzner Cloud, and I wanted to be able to manage
and repair the cluster from within. There are really good projects out there but none of them quite did what I wanted (as of 2020-09-19).
They are all worth checking out, especially if this project doesn't meet your needs.

* [Pharmer](https://github.com/pharmer/pharmer) - loads of features but doesn't support Hetzner Cloud.
* [Hetzner-Kube](https://github.com/xetys/hetzner-kube) - impressive networking setup dates from before Hetzner Cloud had private networks, load balancers and labels. Uses `kubeadm`.
* [K3sup](https://github.com/alexellis/k3sup) - great way to install Rancher's K3s Kubernetes on a cluster but it doesn't provision the cluster or up a firewall.
* [Kube-Hetzner](https://github.com/mysticaltech/kube-hetzner) - uses Terraform to set up K3OS
* [kubernetes-on-hetzner](https://github.com/LWJ/kubernetes-on-hetzner) - uses Terraform
* [hetzner-k3s](https://github.com/vitobotta/hetzner-k3s) - Vito Botta's more recent project - very complete but manages the cluster from outside
