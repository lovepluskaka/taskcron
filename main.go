package main

import (
	"context"
	"github.com/lovepluskaka/taskcron/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type TaskService struct {
}

func (t TaskService) Create(ctx context.Context, in *proto.CreateRequest) (*proto.CreateResponse, error) {
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	s := grpc.NewServer()
	proto.RegisterTaskServiceServer(s, TaskService{})
	s.Serve(lis)
}
