FROM golang:1.17-alpine as build
WORKDIR /go/src
ADD . /go/src
RUN go get -d -v ./...
RUN go build -v -o /go/bin/hetzanetes

FROM alpine
COPY --from=build /go/bin/hetzanetes /hetzanetes
ENTRYPOINT ["/hetzanetes"]