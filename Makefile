buildProto:
	protoc --proto_path=. --proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./app --grpc-gateway_out=./app proto/*.proto

go-build:
	GO111MODULE=on go build -o cmd/api/api ./cmd/api/
	GO111MODULE=on go build -o cmd/gw/gw ./cmd/gw/
