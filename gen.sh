#!/usr/bin/env bash
protoc --proto_path=./protos protos/main.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc --go-grpc_out=./lib/protos
protoc --proto_path=./protos protos/main.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go --go_out=./lib/protos