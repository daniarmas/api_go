#!/usr/bin/env bash
protoc --proto_path=./api api/*.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc --go-grpc_out=./
protoc --proto_path=./api api/*.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go --go_out=./

# protoc --proto_path=./api api/*.proto --plugin=$(go env GOPATH)/bin/protoc-gen-validate --validate_out="lang=go:./"