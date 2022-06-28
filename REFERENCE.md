DOCKER:

* https://hub.docker.com/r/duncanpierce/hetzanetes/tags

CONFIG:

* https://blog.atomist.com/kubernetes-apply-replace-patch/
* https://evancordell.com/posts/kube-apis-crds/
* https://kubernetes.io/docs/tasks/run-application/access-api-from-pod/
  * /var/run/secrets/kubernetes.io/serviceaccount/token
  * /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
  * /var/run/secrets/kubernetes.io/serviceaccount/namespace
* https://docs.cilium.io/en/v1.11/gettingstarted/k3s/
* https://docs.cilium.io/en/stable/gettingstarted/host-firewall/#host-firewall
* [Cilium system requirements](https://docs.cilium.io/en/stable/operations/system_requirements/#mounted-ebpf-filesystem)
* [Cilium help](https://docs.cilium.io/en/stable/gettinghelp/)  
* kubectl get secret hetzanetes -n kube-system -o yaml
* https://rancher.com/docs/k3s/latest/en/installation/install-options/server-config/
* https://rancher.com/docs/k3s/latest/en/installation/install-options/
  * `INSTALL_K3S_VERSION`
  * `INSTALL_K3S_CHANNEL` default `stable`
  * Latest versions: https://update.k3s.io/v1-release/channels
* https://rancher.com/docs/k3s/latest/en/installation/ha/
* https://rancher.com/docs/k3s/latest/en/installation/ha-embedded/
* https://rancher.com/docs/k3s/latest/en/installation/installation-requirements/
  * Server sizes and firewall needs
* [client-side API LB](https://www.youtube.com/watch?app=desktop&v=1Fet0qZdQrM)
* https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/#add-a-label-to-a-node
* https://pkg.go.dev/github.com/hetznercloud/hcloud-cloud-controller-manager/internal/annotation
* [Kubernetes REST API docs](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/)
* [Example of strategic patch using REST API](https://stackoverflow.com/questions/71874714/patch-through-kuberentes-rest-api)

CLUSTER MANAGEMENT:

* https://github.com/Masterminds/semver
* https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler
  * https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-the-key-best-practices-for-running-cluster-autoscaler
  * https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-expanders
* `kubectl cluster-info dump`
* `kubectl get pods --all-namespaces -o custom-columns=NAME:.metadata.name,CONTROLLER:.metadata.ownerReferences[].kind,NAMESPACE:.metadata.namespace`
* `kubectl get pods --all-namespaces -o wide --field-selector spec.nodeName=<node-name>`
  * equiv of API query `/api/v1/namespaces/<namespace>/pods?fieldSelector=spec.nodeName%3D<node-name>`
  * then walk through, ignore any with `metadata.ownerReferences[].{kind=DaemonSet,apiVersion=apps/v1}`
  * evict the remaining pods: `POST /api/v1/namespaces/{namespace}/pods/{name}/eviction` (https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#create-eviction-pod-v1-core)
  * cordon a node with:
* [Well known taints and annotations](https://kubernetes.io/docs/reference/labels-annotations-taints/)
  * Taint for control plane nodes is `node-role.kubernetes.io/master:NoSchedule`
  
```yaml
spec:
  taints:
  - effect: NoSchedule
    key: node.kubernetes.io/unschedulable
  unschedulable: true
```
  
LOGS/STATUS:

* /var/log/syslog
* KUBECONFIG=/etc/rancher/k3s/k3s.yaml cilium status

CRD:

* https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/
* https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
* https://developers.redhat.com/blog/2020/12/16/create-a-kubernetes-operator-in-golang-to-automatically-manage-a-simple-stateful-application#build_and_initialize_the_kubernetes_operator
* https://book.kubebuilder.io/
* https://kubernetes.io/docs/reference/using-api/api-concepts/#server-side-apply

HETZNER CCM:

* https://github.com/hetznercloud/hcloud-cloud-controller-manager/blob/master/internal/hcops/load_balancer.go
* https://pkg.go.dev/github.com/hetznercloud/hcloud-cloud-controller-manager/internal/annotation#Name

UPGRADES

* https://wiki.debian.org/UnattendedUpgrades
  * `APT::Periodic::Update-Package-Lists "1";` in `/etc/apt/apt.conf.d/02periodic`
* https://www.linuxcapable.com/how-to-setup-configure-unattended-upgrades-on-ubuntu-20-04/

UNINSTALL:

* /usr/local/bin/k3s-agent-uninstall.sh or /usr/local/bin/k3s-uninstall.sh

EVICTION:

* https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/
* https://stackoverflow.com/questions/57189208/what-are-the-api-involved-during-kubectl-cordon-and-drain-command