package errdetails

import "google.golang.org/grpc/metadata"

type Detail struct {
	Code int32 `json:"code"`
	Name string `json:"name"`
	Message string `json:"message"`
}

type errorBody struct {
	Message string `json:”message"`
	GrpcCode int32 `json:"grpcCode"`
	Details []Detail `json:"details"`
}

// ServerMetadata consists of metadata sent from gRPC server.
type ServerMetadata struct {
	HeaderMD  metadata.MD
	TrailerMD metadata.MD
}