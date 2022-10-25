# #!/usr/bin/env bash
# export PATH=$PATH:$(go env GOPATH)/bin/
# protoc --proto_path=./proto proto/*.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc --go-grpc_out=./
# protoc --proto_path=./proto proto/*.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go --go_out=./

# # grpc-gateway
# # protoc --proto_path=./proto proto/*.proto --grpc-gateway_out ./pkg --grpc-gateway_opt logtostderr=true,allow_delete_body=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true

# # protoc-gen-validate
# # protoc --proto_path=./proto proto/*.proto --go_out=./ --validate_out="lang=go:./"


#!/usr/bin/env bash
export PATH=$PATH:$(go env GOPATH)/bin/
protoc --proto_path=./proto proto/main.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc --go-grpc_out=./
protoc --proto_path=./proto proto/main.proto --plugin=$(go env GOPATH)/bin/protoc-gen-go --go_out=./

# grpc-gateway
# protoc --proto_path=./proto proto/main.proto --grpc-gateway_out ./pkg/grpc --grpc-gateway_opt logtostderr=true,allow_delete_body=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true