package errdetails

type ErrDetail struct {
	Code int32 `json:"code"`
	Name string `json:"name"`
	Message string `json:"message"`
}

type ErrorBody struct {
	Message string      `json:”message"`
	GrpcCode int32      `json:"grpcCode"`
	Details []ErrDetail `json:"details"`
}

const ErrorDetailKey = "error-detail" // converted to lowercase in setTrailer