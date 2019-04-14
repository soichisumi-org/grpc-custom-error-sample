package errdetails

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"strings"
)

func AddErrorDetail(ctx context.Context, detail ErrorDetail) error {
	bytes, err := json.Marshal(detail)
	if err != nil {
		return err
	}

	md := metadata.Pairs(ErrorDetailKey, string(bytes))
	fmt.Printf("added error detail: %+v\n", detail)
	return grpc.SetTrailer(ctx, md)
}

func GetErrorDetails(md runtime.ServerMetadata) ([]ErrorDetail, error) {
	details := make([]ErrorDetail, 0)
	for k, vs := range md.TrailerMD {
		fmt.Printf("k: %+v, vs: %+v\n", k, vs)
		if !strings.Contains(k, ErrorDetailKey) {
			continue
		}

		for _, v := range vs {
			fmt.Printf("v: %+v\n", v)
			var detail ErrorDetail
			if err := json.Unmarshal([]byte(v), &detail); err != nil {
				return nil, err
			}
			details = append(details, detail)
		}
	}
	return details, nil
}

// CustomHTTPError
// ref: https://mycodesmells.com/post/grpc-gateway-error-handler
// ref: https://suzan2go.hatenablog.com/entry/2017/11/14/005350
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"error": "failed to marshal error message"}`
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	body := &ErrorBody{
		Message: s.Message(),
		GrpcCode: int32(s.Code()),
	}

	// set error details to body
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Errorf("Failed to extract ServerMetadata from context")
	}
	details, err := GetErrorDetails(md)
	if err != nil {
		grpclog.Errorf("Failed to get ErrorDetails from metadata")
	}
	body.Details = details

	// marshal body
	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Errorf("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Errorf("Failed to write response: %v", err)
		}
		return
	}

	// convert grpc code to http code
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Errorf("Failed to write response: %v", err)
	}
}
