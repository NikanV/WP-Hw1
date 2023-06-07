package main

import (
	"context"
	"strconv"

	pb "auth-server/auth"
	pb2 "biz-server/biz"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"flag"
)

func makeAuthenticatorClient() (pb.AuthenticatorClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(*authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to dial authenticator-server! " + err.Error())
	}
	return pb.NewAuthenticatorClient(conn), conn
}

func makeBizServiceClient() (pb2.BizServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(*bizServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to dial authenticator-server! " + err.Error())
	}
	return pb2.NewBizServiceClient(conn), conn
}

func reqPQHandler(c *gin.Context) {
	message_id, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	if err != nil {
		panic("Wrong message_id format! " + err.Error())
	} else if message_id%2 != 0 || message_id <= 0 {
		panic("Wrong message_id format! Should be even and greater than zero!")
	}
	nonce := c.Query("nonce")
	if len(nonce) != 20 {
		panic("Wrong nonce format! Should be exactly 20 characters long!")
	}

	client, conn := makeAuthenticatorClient()
	defer conn.Close()
	request := pb.PQRequest{
		Nonce:     nonce,
		MessageId: message_id,
	}

	response, err := client.RequestPQ(context.Background(), &request)
	if err != nil {
		panic("Failed to request PQ! " + err.Error())
	}

	c.JSON(200, gin.H{
		"nonce":        response.Nonce,
		"server_nonce": response.ServerNonce,
		"message_id":   response.MessageId,
		"p":            response.P,
		"g":            response.G,
	})
}

func reqDHParamsHandler(c *gin.Context) {
	message_id, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	if err != nil {
		panic("Wrong message_id format! " + err.Error())
	} else if message_id%2 != 0 || message_id <= 0 {
		panic("Wrong message_id format! Should be even and greater than zero!")
	}
	a, err := strconv.ParseInt(c.Query("a"), 10, 64)
	if err != nil {
		panic("Wrong A format! " + err.Error())
	}
	nonce := c.Query("nonce")
	server_nonce := c.Query("server_nonce")
	if len(nonce) != 20 || len(server_nonce) != 20 {
		panic("Wrong nonce or server_nonce format! Should be exactly 20 characters long!")
	}

	client, conn := makeAuthenticatorClient()
	defer conn.Close()
	request := pb.DHRequest{
		Nonce:       nonce,
		ServerNonce: server_nonce,
		MessageId:   message_id,
		A:           a,
	}
	response, err := client.RequestDHParams(context.Background(), &request)
	if err != nil {
		panic("Failed to request DHParams! " + err.Error())
	}

	c.JSON(200, gin.H{
		"nonce":        response.Nonce,
		"server_nonce": response.ServerNonce,
		"message_id":   response.MessageId,
		"b":            response.B,
	})
}

func getUsersHandler(c *gin.Context) {
	message_id, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	if err != nil {
		panic("Wrong message_id format! " + err.Error())
	} else if message_id%2 != 0 || message_id <= 0 {
		panic("Wrong message_id format! Should be even and greater than zero!")
	}
	user_id, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		panic("Wrong user_id format! " + err.Error())
	}
	auth_key := c.Query("auth_key")
	client, conn := makeBizServiceClient()
	defer conn.Close()
	request := pb2.GetUsersRequest{
		UserId:    user_id,
		AuthKey:   auth_key,
		MessageId: message_id,
	}
	response, err := client.GetUsers(context.Background(), &request)
	if err != nil {
		panic("Failed to get users! " + err.Error())
	}

	c.JSON(200, gin.H{
		"users":      response.Users,
		"message_id": response.MessageId,
	})
}

func getUsersInjectionHandler(c *gin.Context) {
	message_id, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	if err != nil {
		panic("Wrong message_id format! " + err.Error())
	} else if message_id%2 != 0 || message_id <= 0 {
		panic("Wrong message_id format! Should be even and greater than zero!")
	}
	user_id := c.Query("user_id")
	auth_key := c.Query("auth_key")
	client, conn := makeBizServiceClient()
	defer conn.Close()
	request := pb2.GetUsersWithSQLRequest{
		UserId:    user_id,
		AuthKey:   auth_key,
		MessageId: message_id,
	}
	response, err := client.GetUsersWithSQL(context.Background(), &request)
	if err != nil {
		panic("Failed to get users with SQL injection! " + err.Error())
	}

	c.JSON(200, gin.H{
		"users":      response.Users,
		"message_id": response.MessageId,
	})
}

var (
	authServerAddr = flag.String("authAddr", "localhost:5052", "this is the auth server address")
	bizServerAddr  = flag.String("bizAddr", "localhost:5062", "this is the biz server address")
)

func main() {
	flag.Parse()

	r := gin.Default()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Message": "The gateway-server is up",
		})
	})

	r.GET("/auth/reqpq", reqPQHandler)
	r.GET("/auth/reqdh", reqDHParamsHandler)
	r.GET("/biz/getusers", getUsersHandler)
	r.GET("/biz/getusersinjection", getUsersInjectionHandler)

	err := r.Run(":6443")
	if err != nil {
		panic("Failed to start Gin server on port 6443! " + err.Error())
	}
}
