package app

import (
	"context"
	"github.com/soichisumi-sandbox/grpc-custom-error-sample/app/errdetails"
	"github.com/soichisumi-sandbox/grpc-custom-error-sample/app/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer() Server {return Server{}}

type Server struct{}

func (s Server) GetData(ctx context.Context, req *proto.GetDataRequest) (*proto.GetDataResponse, error) {
	if !req.Success {
		_ = errdetails.AddErrorDetail(ctx, errdetails.ErrDetail{
			Code:    1000,
			Name:    "Invalid Request",
			Message: "success is false",
		})
		return &proto.GetDataResponse{}, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &proto.GetDataResponse{
		Str: "result",
		Num: 5,
	}, nil
}
