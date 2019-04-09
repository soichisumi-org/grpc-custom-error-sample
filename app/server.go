package app

import (
	"context"
	"github.com/soichisumi/customErrResponse/app/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const ErrorDetailKeyPrefix = "error-detail-" // converted to lowercase in setTrailer

func NewServer() Server {return Server{}}

type Server struct{}

func (s Server) GetData(ctx context.Context, req *proto.GetDataRequest) (*proto.GetDataResponse, error) {
	if !req.Success {
		md := metadata.Pairs(
			ErrorDetailKeyPrefix + "name", "InvalidParameter",
			ErrorDetailKeyPrefix + "code", "1001",
			ErrorDetailKeyPrefix + "message", "request parameter is invalid")
		_ = grpc.SetTrailer(ctx, md)
		return &proto.GetDataResponse{}, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &proto.GetDataResponse{
		Str: "result",
		Num: 5,
	}, nil
}
