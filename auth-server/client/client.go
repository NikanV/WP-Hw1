package main

import (
	pb "auth-server/auth"
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	tools "tools"
)

func makePQRequest(client pb.AuthenticatorClient) {
	response, err := client.RequestPQ(context.Background(), &pb.PQRequest{Nonce: tools.RandomString(20), MessageId: 4})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Println(response)
}

func makeDHRequest(client pb.AuthenticatorClient) {
	response, err := client.RequestDHParams(context.Background(), &pb.DHRequest{Nonce: "pp111111111111111111", ServerNonce: "tt111111111111111111", MessageId: 6, A: 2})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println(response)
}

var (
	serverAddr = flag.String("addr", "localhost:5052", "this is the server address")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthenticatorClient(conn)

	makePQRequest(client)
	makeDHRequest(client)
}
