# grpc-custom-error-sample

sample implementation of grpc server with grpc-gateway which returns custom error response

## usage

### build

* `make go-build`
* `make buildProto`

### run

* launch `cmd/api/api` and `cmd/gw/gw`

### request example

* server returns response with success
  * `curl localhost:8080/data?success=true` 
* server returns response with detailed error 
  * `curl localhost:8080/data?success=false`
  
## repository layout

* cmd
  * main applications
* proto
  * proto files
* app
  * internal packages