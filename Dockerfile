# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.10.1 AS builder
RUN go version

COPY . /go/src/github.com/Timothylock/inventory-management/
WORKDIR /go/src/github.com/Timothylock/inventory-management/
RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

# Stage 2 (to create a downsized "container executable", ~7MB)
FROM alpine:3.7
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY ./frontend/ /frontend/
COPY --from=builder /go/src/github.com/Timothylock/inventory-management/app .

EXPOSE 9090
ENTRYPOINT ["./app"]