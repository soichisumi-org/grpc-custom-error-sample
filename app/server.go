package app

import (
	"context"
	"github.com/soichisumi/customErrResponse/app/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer() Server {return Server{}}

type Server struct{}

func (s Server) GetData(ctx context.Context, req *proto.GetDataRequest) (*proto.GetDataResponse, error) {
	if !req.Success {
		return &proto.GetDataResponse{}, status.Error(codes.InvalidArgument, "invalid request")
	}
	return &proto.GetDataResponse{
		Str: "result",
		Num: 5,
	}, nil
}
