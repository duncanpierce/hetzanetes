DOCKER:

* https://hub.docker.com/r/duncanpierce/hetzanetes/tags

CONFIG:

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

CRD:

* https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/
* https://www.martin-helmich.de/en/blog/kubernetes-crd-client.html
* https://developers.redhat.com/blog/2020/12/16/create-a-kubernetes-operator-in-golang-to-automatically-manage-a-simple-stateful-application#build_and_initialize_the_kubernetes_operator
* https://book.kubebuilder.io/

HETZNER CCM:

* https://github.com/hetznercloud/hcloud-cloud-controller-manager/blob/master/internal/hcops/load_balancer.go
* https://pkg.go.dev/github.com/hetznercloud/hcloud-cloud-controller-manager/internal/annotation#Name

UPGRADES

* https://wiki.debian.org/UnattendedUpgrades
  * `APT::Periodic::Update-Package-Lists "1";` in `/etc/apt/apt.conf.d/02periodic`
* https://www.linuxcapable.com/how-to-setup-configure-unattended-upgrades-on-ubuntu-20-04/ 