proto:
	protoc --proto_path=. --proto_path=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:./app --grpc-gateway_out=./app proto/*.proto
