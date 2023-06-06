package main

import (
	"context"

	pb "auth-server/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"flag"
)

func reqPQHandler(c *gin.Context) {
	conn, err := grpc.Dial(*authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("fail to dial: " + err.Error())
	}
	defer conn.Close()
	client := pb.NewAuthenticatorClient(conn)

	response, err := client.RequestPQ(context.Background(), &pb.PQRequest{Nonce: "client_nonce", MessageId: 4})
	if err != nil {
		panic("failed to authenticate: " + err.Error())
	}

	c.JSON(200, gin.H{
		"nonce":        response.Nonce,
		"server_nonce": response.ServerNonce,
		"message_id":   response.MessageId,
		"g":            response.G,
	})
}

func reqDHParamsHandler(c *gin.Context) {
	conn, err := grpc.Dial(*authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("fail to dial: " + err.Error())
	}
	defer conn.Close()
	client := pb.NewAuthenticatorClient(conn)

	response, err := client.RequestDHParams(context.Background(), &pb.DHRequest{Nonce: "pp", ServerNonce: "tt", MessageId: 6, A: 2})
	if err != nil {
		panic("failed to send key: " + err.Error())
	}

	c.JSON(200, gin.H{
		"nonce":        response.Nonce,
		"server_nonce": response.ServerNonce,
		"message_id":   response.MessageId,
		"b":            response.B,
	})
}

var (
	authServerAddr = flag.String("addr", "localhost:5052", "this is the server address")
)

func main() {
	flag.Parse()

	r := gin.Default()
	r.GET("/auth/reqpq", reqPQHandler)
	r.GET("/auth/reqdh", reqDHParamsHandler)

	err := r.Run(":6443")
	if err != nil {
		panic("[Error] failed to start Gin server due to: " + err.Error())
	}
}
