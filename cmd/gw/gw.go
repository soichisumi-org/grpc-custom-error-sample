package main

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soichisumi/customErrResponse/app"
	"github.com/soichisumi/customErrResponse/app/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"strings"
)

const grpcEndpoint = ":3000"
const httpEndpoint = ":8080"

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

type applicationMetadataKey struct{}


// NewServerMetadataContext creates a new context with ServerMetadata
func NewServerMetadataContext(ctx context.Context, md ServerMetadata) context.Context {
	return context.WithValue(ctx, applicationMetadataKey{}, md)
}

// ServerMetadataFromContext returns the ServerMetadata in ctx
func ServerMetadataFromContext(ctx context.Context) (md ServerMetadata, ok bool) {
	md, ok = ctx.Value(applicationMetadataKey{}).(ServerMetadata)
	return
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

	body := &errorBody{
		Message: s.Message(),
		GrpcCode: int32(s.Code()),
	}

	// retrieve server metadata
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Infof("Failed to extract ServerMetadata from context")
	}
	for k, vs := range md.TrailerMD {
		fmt.Printf("k: %+v, vs: %+v\n", k, vs)
		if !strings.Contains(k, app.ErrorDetailKeyPrefix) {
			continue
		}
		for _, v := range vs {
			fmt.Printf("v: %+v\n", strings.ReplaceAll(v, app.ErrorDetailKeyPrefix, ""))
		}
	}

	// marshal body
	buf, merr := marshaler.Marshal(body)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", body, merr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	// convert grpc code to http code
	st := runtime.HTTPStatusFromCode(s.Code())
	w.WriteHeader(st)
	if _, err := w.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}


func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(grpcEndpoint, dialOpts...) // grpcサーバのエンドポイント
	if err != nil {
		return nil, err
	}

	err = proto.RegisterServerHandler(ctx, mux, conn)
	if err != nil {
		return nil, err
	}

	corsMux := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "HEAD", "PUT", "DELETE"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "content-type", "authorization"}),
	)(mux)

	return corsMux, nil
}

func Run(address string, opts ...runtime.ServeMuxOption) error {
	// apply custom http error
	runtime.HTTPError = CustomHTTPError

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}

	return http.ListenAndServe(address, gw)
}

func main() {
	fmt.Printf("http api is running on port: %s...\n", httpEndpoint)
	if err := Run(httpEndpoint); err != nil {
		panic(err)
	}
}
