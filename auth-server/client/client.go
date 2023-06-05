package main

import (
	pb "auth-server/auth"
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func makePQRequest(client pb.Req_PQClient) {
	response, err := client.RequestPQ(context.Background(), &pb.PQRequest{Nonce: "client_nonce11111111", MessageId: 4})
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Println(response)
}

func makeDHRequest(client pb.Req_DH_ParamsClient) {
	response, err := client.RequestDHparams(context.Background(), &pb.DHRequest{Nonce: "pp", ServerNonce: "tt" , MessageId: 6 , A: 2})
	if err != nil {
		log.Fatalf("error: %v", err)
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

	client := pb.NewReq_PQClient(conn)
	client2 := pb.NewReq_DH_ParamsClient(conn)

	makePQRequest(client)
	makeDHRequest(client2)
	
}
