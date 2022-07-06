package gg

import (
	"context"
	"github.com/b1uem0nday/transfer_service/internal/contracts"
	p "github.com/b1uem0nday/transfer_service/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GrpcGateway struct {
	transferClient *contracts.Worker
	ctx            context.Context
	gs             *grpc.Server
	listener       net.Listener
	p.UnimplementedTransferServiceServer
}

func New(ctx context.Context, client *contracts.Worker) *GrpcGateway {
	return &GrpcGateway{
		transferClient:                     client,
		ctx:                                ctx,
		UnimplementedTransferServiceServer: p.UnimplementedTransferServiceServer{},
	}
}

func (gg *GrpcGateway) Connect(port string) (err error) {
	gg.listener, err = net.Listen("tcp", ":"+port)

	if err != nil {
		return err
	}

	var opts []grpc.ServerOption

	gg.gs = grpc.NewServer(opts...)
	p.RegisterTransferServiceServer(gg.gs, gg)
	return nil
}

func (gg *GrpcGateway) Run() (err error) {
	log.Println("listening", gg.listener.Addr())
	go gg.gs.Serve(gg.listener)
	<-gg.ctx.Done()
	return nil
}