package main

import (
	pb "biz-server/biz"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	serverAddr = flag.String("addr", "localhost:5062", "this is the server address")
)

func makeGetUserReq(client pb.Get_UsersClient) {
	response, err := client.GetUsers(context.Background(),
		&pb.Get_Users_Req{UserId: 9649, AuthKey: "authKey", MessageId: 4})
	// todo: check auth key
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}
	log.Println(response)
}

func makeGetUserWSqlInjReq(client pb.Get_UsersClient) {
	response, err := client.Get_User_Sql_Inj(context.Background(),
		&pb.Get_User_Sql_Inj_Req{AuthKey: "authKey", MessageId: 4})
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
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

	client := pb.NewGet_UsersClient(conn)
	//client2 := pb.NewGet_UsersClient(conn)
	//
	makeGetUserReq(client)
	//makeGetUserReq(client2)
}
