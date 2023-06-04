package main

import (
	pb "auth-server/auth"
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sendAsadi(client pb.RepPqClient) {
	response, err := client.Authenticate(context.Background(), &pb.Request{Name: "asadi"})
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}

	log.Println(response)
}

var (
	serverAddr = flag.String("addr", "localhost:8080", "this is the server address")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewRepPqClient(conn)

	sendAsadi(client)
}
