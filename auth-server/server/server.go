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

type reqPqServer struct {
	pb.UnimplementedReqPqServer
}

func (c *reqPqServer) RequestPQ(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	fmt.Println("Got the request")
	return &pb.Response{Nonce: in.GetNonce(), ServerNonce: "server_nonce", MessageId: in.GetMessageId() + 1, P: 23, G: 5}, nil
}

func newServer() *reqPqServer {
	s := &reqPqServer{}
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
	} else {
		fmt.Println("listening on port: ", *port)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterReqPqServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
