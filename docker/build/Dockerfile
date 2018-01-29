FROM alpine:3.7

MAINTAINER own3dh2so4 <own3dh2so4@gmail.com>

RUN apk add --update go-tools git make protobuf && \
    go get -u 'github.com/golang/protobuf/protoc-gen-go' && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway && \
    go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger && \
    go get -u gopkg.in/alecthomas/gometalinter.v2 && gometalinter.v2 --install && \
    go get -u github.com/axw/gocov/gocov && \
    go get -u github.com/AlekSi/gocov-xml && \
    go get -u github.com/jstemmer/go-junit-report