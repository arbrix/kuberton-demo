FROM golang:1.11-alpine as builder
RUN apk add --no-cache ca-certificates git

ENV PROJECT github.com/arbrix/kubertron-demo
WORKDIR /go/src/$PROJECT

# restore dependencies
COPY go.* ./
RUN GO111MODULE=on go mod vendor
COPY . .
RUN go install .

FROM alpine as release
RUN apk add --no-cache ca-certificates \
    busybox-extras net-tools bind-tools
WORKDIR /shop
COPY --from=builder /go/bin/kubertron-demo /shop/server
COPY ./templates ./templates
COPY ./static ./static
EXPOSE 3000
ENTRYPOINT ["/shop/server"]
