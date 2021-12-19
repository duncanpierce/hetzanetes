DOCKER:

* https://hub.docker.com/r/duncanpierce/hetzanetes/tags

CONFIG:

* kubectl get secret hetzanetes -n kube-system -o yaml
* https://rancher.com/docs/k3s/latest/en/installation/install-options/server-config/
* https://rancher.com/docs/k3s/latest/en/installation/ha/
* https://rancher.com/docs/k3s/latest/en/installation/ha-embedded/
* [client-side API LB](https://www.youtube.com/watch?app=desktop&v=1Fet0qZdQrM)
* https://kubernetes.io/docs/tasks/configure-pod-container/assign-pods-nodes/#add-a-label-to-a-node
* https://pkg.go.dev/github.com/hetznercloud/hcloud-cloud-controller-manager/internal/annotation

HETZNER CCM:

* https://github.com/hetznercloud/hcloud-cloud-controller-manager/blob/master/internal/hcops/load_balancer.go
* https://pkg.go.dev/github.com/hetznercloud/hcloud-cloud-controller-manager/internal/annotation#Name
