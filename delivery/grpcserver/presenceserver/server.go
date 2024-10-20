package presenceserver

import (
	"context"
	"fmt"
	"gapp/contract/goproto/presence"
	"gapp/param"
	"gapp/pkg/protobufmapper"
	"gapp/pkg/slice"
	"gapp/service/presenceservice"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	presence.UnimplementedPresenceServiceServer
	svc presenceservice.Service
}

func New(svc presenceservice.Service) Server {
	return Server{
		UnimplementedPresenceServiceServer: presence.UnimplementedPresenceServiceServer{},
		svc:                                svc,
	}
}

func (s Server) GetPresence(ctx context.Context, req *presence.GetPresenceRequest) (*presence.GetPresenceResponse, error) {
	resp, err := s.svc.GetPresence(ctx, param.GetPresenceRequest{
		UserIDs: slice.MapFromUint64ToUint(req.GetUserIds()),
	})

	if err != nil {
		return nil, err
	}

	return protobufmapper.MapGetPresenceResponseToProtobuf(resp), nil
}

func (s Server) Start() {
	// listener := tcp port
	address := fmt.Sprintf(":%d", 8086)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	// grpc server
	grpcServer := grpc.NewServer()
	// pbPresenceserver register into grpc server

	presence.RegisterPresenceServiceServer(grpcServer, &s)
	// server grpcServer by listener

	log.Println("presence grpc server starting on", address)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("couldn't server presence grpc server")
	}

}
