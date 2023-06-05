package main

import (
<<<<<<< HEAD
	pb "biz-server/biz"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type getUserServer struct {
	pb.UnimplementedGet_UsersServer
}

func (c *getUserServer) GetUsers(ctx context.Context,
	in *pb.Get_Users_Req) (*pb.Get_Users_Resp, error) {
	fmt.Println("Get user request received")
	//users := postgres get user with userid s
	// if userid is nil return top 100 users
	// check if user's auth key is similar to input's auth key
	//return &pb.Get_Users_Response{Users: users, MessageId: in.GetMessageId() + 1}, nil
	return nil, nil
}

func newGetUserServer() *getUserServer {
	return &getUserServer{}
}

////////////////////////////////second service

func (c *getUserServer) GetUserWSqlInj(ctx context.Context,
	in *pb.Get_User_Sql_Inj_Req) (*pb.Get_Users_Resp, error) {
	fmt.Println("Get user with sql injection request received")
	//return &pb.Get_Users_Response{Users:, MessageId: in.GetMessageId() + 1},nil
	return nil, nil
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
	pb.RegisterGet_UsersServer(grpcServer, newGetUserServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
=======
	"fmt"
)

func main() {
	fmt.Println("hello world!")
>>>>>>> nikan
}
