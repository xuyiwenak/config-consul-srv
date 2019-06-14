FROM golang:1.12.4-alpine3.9 as builder

WORKDIR /go/src/consul-config-push

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o consul-config-push loader.go main.go

FROM alpine:3.9 as prod

WORKDIR /root/consul-config-push

COPY --from=0 /go/src/consul-config-push  .

CMD ["./consul-config-push"]