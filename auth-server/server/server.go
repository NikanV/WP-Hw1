package main

import (
	pb "auth-server/auth"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type repPqServer struct {
	pb.UnimplementedRepPqServer
}

func (a *repPqServer) Authenticate(context.Context, *pb.Request) (*pb.Reply, error) {
	fmt.Println("omad")
	return &pb.Reply{Response: "daddy"}, nil
}

func newServer() *repPqServer {
	s := &repPqServer{}
	return s
}

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Listening on port: ", *port)

	grpcServer := grpc.NewServer()
	pb.RegisterRepPqServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
