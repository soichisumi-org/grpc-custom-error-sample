package main

import (
	"fmt"
	"github.com/soichisumi-sandbox/grpc-custom-error-sample/app"
	"github.com/soichisumi-sandbox/grpc-custom-error-sample/app/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	grpcPort = ":3000"
)

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %+v\n", err)
	}

	grpcServer := grpc.NewServer()
	server := app.NewServer()
	proto.RegisterServerServer(grpcServer, server)
	reflection.Register(grpcServer)
	fmt.Printf("api is running on port: %s\n", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
