FROM alpine
COPY bin/hetzanetes /hetzanetes
ENTRYPOINT ["/hetzanetes"]
LABEL org.opencontainers.image.description="A simple way to set up and manage Kubernetes clusters on Hetzner Cloud."