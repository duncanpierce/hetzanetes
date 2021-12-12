FROM alpine
COPY bin/hetzanetes /hetzanetes
ENTRYPOINT ["/hetzanetes"]