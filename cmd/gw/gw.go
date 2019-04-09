package main

import (
	"context"
	"fmt"
	"github.com/soichisumi/customErrResponse/app/proto"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"net/http"
)

const grpcEndpoint = ":3000"
const httpEndpoint = ":8080"

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
