FROM alpine
COPY bin/hetzanetes /hetzanetes
ENTRYPOINT ["/hetzanetes"]
LABEL org.opencontainers.image.description="Create self-managing Rancher K3s Kubernetes clusters on Hetzner Cloud."