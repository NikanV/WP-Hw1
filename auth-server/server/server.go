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

type authServer struct {
	pb.UnimplementedAuthenticatorServer
}

func (c *authServer) RequestPQ(ctx context.Context, in *pb.PQRequest) (*pb.PQResponse, error) {
	fmt.Println("Got the request")
	return &pb.PQResponse{Nonce: in.GetNonce(), ServerNonce: "server_nonce", MessageId: in.GetMessageId() + 1, P: 23, G: 5}, nil
}

func (c *authServer) RequestDHParams(ctx context.Context, in *pb.DHRequest) (*pb.DHResponse, error) {
	fmt.Println("got the public key")
	return &pb.DHResponse{Nonce: in.GetNonce(), ServerNonce: in.GetServerNonce(), MessageId: in.GetMessageId() + 1, B: 22}, nil
}

func newAuthServer() *authServer {
	s := &authServer{}
	return s
}

var (
	port = flag.Int("port", 5052, "The server port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("listening on port: ", *port)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, newAuthServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
