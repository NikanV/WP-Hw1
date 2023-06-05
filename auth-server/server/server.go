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

type reqPqServer struct {
	pb.UnimplementedReq_PQServer
}

func (c *reqPqServer) RequestPQ(ctx context.Context, in *pb.PQRequest) (*pb.PQResponse, error) {
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

func newPQServer() *reqPqServer {
	s := &reqPqServer{}
	return s
}

////////////////////////////////DH service

type  reqDHParamsServer struct {
	pb.UnimplementedReq_DH_ParamsServer
}

func (c *reqDHParamsServer) RequestDHparams(ctx context.Context, in *pb.DHRequest) (*pb.DHResponse, error) {
	fmt.Println("got the public key")
	return &pb.DHResponse{Nonce: in.GetNonce() , ServerNonce: in.GetServerNonce() , MessageId: in.GetMessageId() + 1 , B: 22} , nil
}

func newDHServer() *reqDHParamsServer {
	s := &reqDHParamsServer{}
	return s
}

///////////////////////////////////


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
	pb.RegisterReq_PQServer(grpcServer, newPQServer())
	pb.RegisterReq_DH_ParamsServer(grpcServer, newDHServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
