package main

import (
	pb "auth-server/auth"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	tools "tools"

	"google.golang.org/grpc"
)

type authServer struct {
	pb.UnimplementedAuthenticatorServer
}

func (c *authServer) RequestPQ(ctx context.Context, in *pb.PQRequest) (*pb.PQResponse, error) {
	message_id := in.GetMessageId()
	nonce := in.GetNonce()
	if len(nonce) != 20 {
		return nil, errors.New("PQRequest : wrong nonce format")
	} else if message_id%2 != 0 || message_id <= 0 {
		return nil, errors.New("PQRequest : wrong message_id format")
	}
	fmt.Println(in)

	prime := tools.Random_Prime()
	response := pb.PQResponse{
		Nonce:       nonce,
		ServerNonce: tools.RandomString(20),
		MessageId:   message_id + 1,
		P:           int64(prime),
		G:           int64(tools.FindPrimitive(prime)),
	}
	return &response, nil
}

func (c *authServer) RequestDHParams(ctx context.Context, in *pb.DHRequest) (*pb.DHResponse, error) {
	message_id := in.GetMessageId()
	nonce := in.GetNonce()
	server_nonce := in.GetServerNonce()
	a := in.GetA()
	if len(nonce) != 20 || len(server_nonce) != 20 {
		return nil, errors.New("DHRequest : wrong nonce or server_nonce format")
	} else if message_id%2 != 0 || message_id <= 0 {
		return nil, errors.New("DHRequest : wrong message_id format")
	}
	fmt.Println(in)

	private_key := tools.RandomNumber(50)
	p, g := 23, 5
	auth_key := int64(p) % int64(math.Pow(float64(a), float64(private_key)))
	fmt.Println(auth_key)

	response := pb.DHResponse{
		Nonce:       nonce,
		ServerNonce: server_nonce,
		MessageId:   message_id,
		B:           int64(math.Pow(1, float64(private_key))) % int64(g),
	}
	return &response, nil
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
		log.Fatalf("Failed to listen! %v", err)
	} else {
		fmt.Println("Listening on port: ", *port)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, newAuthServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve! %v", err)
	}
}
