package main

import (
	pb "biz-server/biz"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func makeGetUserRequest(client pb.Get_UsersClient) {
	response, err := client.GetUsers(context.Background(),
		&pb.Get_Users_Req{AuthKey: "authKey", MessageId: 4})
	// check userid
	// check authKey
	if err != nil {
		log.Fatalf("failed to authenticate: %v", err)
	}
	log.Println(response)
}

func makeGetUserWithSqlInjRequest(client pb.Get_UsersClient) {
	response, err := client.Get_User_Sql_Inj(context.Background(),
		&pb.Get_User_Sql_Inj_Req{AuthKey: "authKey", MessageId: 4})
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

	client := pb.NewGet_UsersClient(conn)

	makeGetUserRequest(client)

}
