package main

import (
	pb "biz-server/biz"
	"context"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var (
	port  = flag.Int("port", 5062, "The server port")
	db    *gorm.DB
	err   error
	users = []*pb.USERS{
		{Name: "nikan", Family: "vsi", Id: 5303, Age: 19, Sex: "male", CreatedAt: time.Now().UTC().String()},
		{Name: "nima", Family: "enigma", Id: 5263, Age: 19, Sex: "male", CreatedAt: time.Now().UTC().String()},
		{Name: "alireza", Family: "kaz", Id: 9649, Age: 19, Sex: "male", CreatedAt: time.Now().UTC().String()},
	}
)

const (
	db_host = "localhost"
	db_port = 5432
	db_user = "postgres"
	db_pass = "postgres"
	db_name = "postgres"
)

//type USERS struct {
//	gorm.Model
//	name      string `json:"name"`
//	family    string `jason:"family"`
//	id        int    `json:"id"`
//	age       int    `json:"age"`
//	sex       string `json:"sex"`
//	createdAt string `json:"createdAt"`
//}

type getUserServer struct {
	pb.UnimplementedGet_UsersServer
}

func (c *getUserServer) GetUsers(ctx context.Context,
	in *pb.Get_Users_Req) (*pb.Get_Users_Resp, error) {
	fmt.Println("Get user request received")
	var usr []*pb.USERS
	if &in.UserId != nil {
		var current *pb.USERS
		db.First(&current, in.UserId)
		usr = append(usr, current)
	} else {
		var allUsers []*pb.USERS
		db.Find(&allUsers)
		minOf := 99
		if len(allUsers) < 99 {
			minOf = len(allUsers)
		}
		usr = allUsers[0:minOf]
	}
	return &pb.Get_Users_Resp{Users: usr, MessageId: in.GetMessageId() + 1}, nil
}

func newGetUserServer() *getUserServer {
	return &getUserServer{}
}

func (c *getUserServer) GetUserWSqlInj(ctx context.Context,
	in *pb.Get_User_Sql_Inj_Req) (*pb.Get_Users_Resp, error) {
	fmt.Println("Get user with sql injection request received")
	//return &pb.Get_Users_Response{Users:, MessageId: in.GetMessageId() + 1},nil
	return nil, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println("failed to connect to database")
		log.Fatal(err)
	}
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_pass, db_name)

	db, err := gorm.Open("postgres", psqlconn)

	checkError(err)

	fmt.Println("connected to database!")

	defer db.Close()

	db.AutoMigrate(pb.USERS{})

	for i := range users {
		db.Create(&users[i])
	}

	//var allUsers []*pb.USERS
	//db.Find(&allUsers)
	//fmt.Println(allUsers)

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
}
