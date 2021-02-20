FROM golang:1.16 as build
WORKDIR /go/src
ADD . /go/src
RUN go get -d -v ./...
RUN go build -v -o /go/bin/hetzanetes

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/hetzanetes /hetzanetes
ENTRYPOINT ["/hetzanetes"]