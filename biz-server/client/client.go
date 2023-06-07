package main

import (
	pb "biz-server/biz"
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:5062", "this is the server address")
)

func makeGetUsersRequest(client pb.BizServiceClient) {
	response, err := client.GetUsers(context.Background(), &pb.GetUsersRequest{UserId: 9649, AuthKey: "authKey", MessageId: 4})
	if err != nil {
		log.Fatalf("Failed to get users! %v", err)
	}

	log.Println(response)
}

func makeGetUsersWithSQLRequest(client pb.BizServiceClient) {
	response, err := client.GetUsersWithSQL(context.Background(), &pb.GetUsersWithSQLRequest{UserId: "9649", AuthKey: "authKey", MessageId: 4})
	if err != nil {
		log.Fatalf("Failed to get users with SQL: %v", err)
	}
	log.Println(response)
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	fmt.Println("connected to the server")
	defer conn.Close()

	client := pb.NewBizServiceClient(conn)
	makeGetUsersRequest(client)
	makeGetUsersWithSQLRequest(client)
}
