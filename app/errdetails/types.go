package errdetails

type ErrorDetail struct {
	Code    int32  `json:"code"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ErrorBody struct {
	Message  string        `json:”message"`
	GrpcCode int32         `json:"grpcCode"`
	Details  []ErrorDetail `json:"details"`
}

const ErrorDetailKey = "error-detail" // converted to lowercase in setTrailer

// Error definitions
// 1000 .. Validation error
// 2000 .. Application error
// 3000 .. Security error
// 4000 .. System error
// 5000 .. Unknown

var (
	// 2000
	SuccessIsFalse = ErrorDetail{Name: "ErrorDetail", Code: 2000, Message: "success param is false"}
)
