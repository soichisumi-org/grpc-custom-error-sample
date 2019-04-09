package app

import (
	"context"
	"github.com/soichisumi/customErrResponse/app/proto"
)

type Server struct{}

func (Server) GetData(context.Context, *proto.GetDataRequest) (*proto.GetDataResponse, error) {
	return
}

