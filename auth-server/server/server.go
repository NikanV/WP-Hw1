package main

import (
	pb "auth-server/auth"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"errors"
	"google.golang.org/grpc"
	tools "tools"
)

type authServer struct {
	pb.UnimplementedAuthenticatorServer
}

func (c *authServer) RequestPQ(ctx context.Context, in *pb.PQRequest) (*pb.PQResponse, error) {
	if(len(in.GetNonce()) != 20){
		fmt.Println("PQRequest : wrong nonce format")
		return nil , errors.New("PQRequest : wrong nonce format")
	}
	if(in.GetMessageId() <= 0 || in.GetMessageId() % 2 == 1){
		fmt.Println("PQRequest : wrong message_id format")
		return nil , errors.New("PQRequest : wrong message_id format")
	}
	fmt.Println(in)
	return &pb.PQResponse{Nonce: in.GetNonce(), ServerNonce: tools.RandomString(20) , MessageId: in.GetMessageId() + 1, P: 23, G: 5}, nil
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
