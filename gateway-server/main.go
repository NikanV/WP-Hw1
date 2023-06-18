package main

import (
	"context"
	"strconv"
	"time"

	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/files" // swagger embed files
	_ "WP-Hw1/docs"

	pb "WP-Hw1/proto"

	rl "github.com/JGLTechnologies/gin-rate-limit"
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

func makeBizServiceClient() (pb.BizServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(*bizServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("Failed to dial authenticator-server! " + err.Error())
	}
	return pb.NewBizServiceClient(conn), conn
}


// reqpq godoc

// @Summary Request PQ
// @Description first step of registeration which we send user info.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param message_id query int true "The message ID (even and greater than zero)."
// @Param nonce query string true "The nonce (20 characters long)."
// @Success 200 {object} pb.PQResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/reqpq [get]
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

// reqdh godoc

// @Summary Request DH
// @Description second step of registeration which we send auth info and public keys.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param message_id query int true "The message ID (even and greater than zero)."
// @Param nonce query string true "The nonce (20 characters long)."
// @Param server_nonce query string true "The server_nonce (20 characters long)."
// @Param a query int true "public key from client"
// @Success 200 {object} pb.DHResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/reqdh [get]
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

func authCheckHandler(c *gin.Context) bool {
	message_id, err := strconv.ParseInt(c.Query("message_id"), 10, 64)
	if err != nil {
		panic("Wrong message_id format! " + err.Error())
	} else if message_id%2 != 0 || message_id <= 0 {
		panic("Wrong message_id format! Should be even and greater than zero!")
	}
	nonce := c.Query("nonce")
	server_nonce := c.Query("server_nonce")
	if len(nonce) != 20 || len(server_nonce) != 20 {
		panic("Wrong nonce or server_nonce format! Should be exactly 20 characters long!")
	}
	auth_key, err := strconv.ParseInt(c.Query("auth_key"), 10, 64)
	if err != nil {
		panic("Wrong auth_key format! " + err.Error())
	}
	client, conn := makeAuthenticatorClient()
	defer conn.Close()
	request := pb.ACRequest{
		Nonce:       nonce,
		ServerNonce: server_nonce,
		MessageId:   message_id,
		AuthKey:     auth_key,
	}
	response, err := client.AuthCheck(context.Background(), &request)
	if err != nil {
		panic("Failed to AuthCheck! " + err.Error())
	}

	return response.AuthCheck
}

// getUsers godoc

// @Summary get users of database.
// @Description after checking authentication , gets the information that you desire.
// @Tags Biz server
// @Accept json
// @Produce json
// @Param user_id query int true "gets first 100 users if negetive"
// @Param auth_key query int true "auth key"
// @Param message_id query int true "The message ID (even and greater than zero)."
// @Success 200 {object} pb.GetUsersResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /biz/getusers [get]
func getUsersHandler(c *gin.Context) {
	if authCheckHandler(c) {
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
		request := pb.GetUsersRequest{
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
	} else {
		c.JSON(404, gin.H{
			"error message": "Invalid authentication key!",
		})
	}
}

// getUsersByInjection godoc

// @Summary get users of database by injection.
// @Description after checking authentication , gets the information that you desire by injection.
// @Tags Biz server
// @Accept json
// @Produce json
// @Param user_id query string true "gets first 100 users if negetive"
// @Param auth_key query int true "auth key"
// @Param message_id query int true "The message ID (even and greater than zero)."
// @Success 200 {object} pb.GetUsersResponse
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /biz/getusersinjection [get]
func getUsersInjectionHandler(c *gin.Context) {
	if authCheckHandler(c) {
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
		request := pb.GetUsersWithSQLRequest{
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
	} else {
		c.JSON(404, gin.H{
			"error message": "Invalid authentication key!",
		})
	}
}

func ipHandler(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info rl.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

var (
	authServerAddr = flag.String("authAddr", "localhost:5052", "this is the auth server address")
	bizServerAddr  = flag.String("bizAddr", "localhost:5062", "this is the biz server address")
)


// @title WebPrograming homework 1
// @description a service which you can register in and get access to the users database
// @version 1.0

func main() {
	flag.Parse()

	r := gin.Default()

	store := rl.InMemoryStore(&rl.InMemoryOptions{
		Rate:  time.Second,
		Limit: 100,
	})

	limiter := rl.RateLimiter(store, &rl.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      ipHandler,
	})

	r.GET("/test", limiter, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Message": "The gateway-server is up",
		})
	})
	r.GET("/auth/reqpq", limiter, reqPQHandler)
	r.GET("/auth/reqdh", limiter, reqDHParamsHandler)
	r.GET("/biz/getusers", limiter, getUsersHandler)
	r.GET("/biz/getusersinjection", limiter, getUsersInjectionHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":6443")
	if err != nil {
		panic("Failed to start Gin server on port 6443! " + err.Error())
	}
}
