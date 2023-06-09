package main

import (
	pb "auth-server/auth"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"time"
	tools "tools"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type authServer struct {
	pb.UnimplementedAuthenticatorServer
}

func (c *authServer) AuthCheck(ctx context.Context, in *pb.ACRequest) (*pb.ACResponce, error) {
	message_id := in.GetMessageId()
	nonce := in.GetNonce()
	server_nonce := in.GetServerNonce()
	auth_key := in.GetAuthKey()
	auth_check := true
	if len(nonce) != 20 || len(server_nonce) != 20 {
		return nil, errors.New("ACRequest : wrong nonce or server_nonce format")
	} else if message_id%2 != 0 || message_id <= 0 {
		return nil, errors.New("ACRequest : wrong message_id format")
	}
	client := initRedisClient()
	defer client.Close()
	hash := tools.Sha1_gen(nonce+server_nonce)
	server_auth_key_str, err := client.Get(context.Background(), hash).Result()
	if(err == redis.Nil){
		return nil, errors.New("no key found!")
	} 
	if(err != nil){
		return nil, errors.New("wrong key format")
	}
	server_auth_key, _ := strconv.ParseInt(server_auth_key_str , 10 , 64)
	fmt.Println("the server-auth-key is : %d\n" , server_auth_key)
	if(server_auth_key != auth_key){
		auth_check = false
	}
	response := pb.ACResponce{
		MessageId:   message_id + 1,
		AuthCheck:    auth_check,
	}
	return &response, nil
}

func (c *authServer) RequestPQ(ctx context.Context, in *pb.PQRequest) (*pb.PQResponse, error) {
	message_id := in.GetMessageId()
	nonce := in.GetNonce()
	if len(nonce) != 20 {
		return nil, errors.New("PQRequest : wrong nonce format")
	} else if message_id%2 != 0 || message_id <= 0 {
		return nil, errors.New("PQRequest : wrong message_id format")
	}
	prime := tools.Random_Prime()
	p, g := int64(prime), int64(tools.FindPrimitive(prime))
	server_nonce := tools.RandomString(20)

	client := initRedisClient()
	defer client.Close()
	hash := tools.Sha1_gen(nonce+server_nonce)
	fmt.Println(hash)
	err := client.HSet(context.Background(), hash , "p" , p).Err()
	if err != nil {
		return nil, err
	}
	err = client.HSet(context.Background(), hash , "g" , g).Err()
	if err != nil {
		return nil, err
	}
	client.Expire(context.Background() , hash , 2*time.Minute)


	fmt.Println(in)
	response := pb.PQResponse{
		Nonce:       nonce,
		ServerNonce: server_nonce,
		MessageId:   message_id + 1,
		P:           p,
		G:           g,
	}
	return &response, nil
}

func (c *authServer) RequestDHParams(ctx context.Context, in *pb.DHRequest) (*pb.DHResponse, error) {
	message_id := in.GetMessageId()
	nonce := in.GetNonce()
	server_nonce := in.GetServerNonce()
	a := in.GetA()
	if len(nonce) != 20 || len(server_nonce) != 20 {
		return nil, errors.New("DHRequest : wrong nonce or server_nonce format")
	} else if message_id%2 != 0 || message_id <= 0 {
		return nil, errors.New("DHRequest : wrong message_id format")
	}
	private_key := tools.RandomNumber(10)
	fmt.Println(private_key)
	hash := tools.Sha1_gen(nonce+server_nonce)
	fmt.Println(hash)
	client := initRedisClient()
	defer client.Close()
	//data := client.HGetAll(context.Background(), hash).Val()

	p, _ := strconv.ParseInt(client.HGet(context.Background() , hash , "p").Val(), 10, 64)
	g, _ := strconv.ParseInt(client.HGet(context.Background() , hash , "g").Val(), 10, 64)
	fmt.Println(int64(p) , math.Pow(float64(a), float64(private_key)) , float64(a) , float64(private_key))
	auth_key := int64(p) % int64(math.Pow(float64(a), float64(private_key)))
	client.Del(context.Background(), hash)
	err := client.Set(context.Background(), hash, auth_key, 0).Err()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Auth key is %d\n", auth_key)
	fmt.Println(in)
	response := pb.DHResponse{
		Nonce:       nonce,
		ServerNonce: server_nonce,
		MessageId:   message_id + 1,
		B:           int64(math.Pow(float64(g), float64(private_key))) % int64(p),
	}
	return &response, nil
}

func initRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func newAuthServer() *authServer {
	s := &authServer{}
	return s
}

var (
	port = flag.Int("port", 5052, "The server port")
)

func main() {
	flag.Parse()

	initRedisClient()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen! %v", err)
	} else {
		fmt.Println("Listening on port: ", *port)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthenticatorServer(grpcServer, newAuthServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve! %v", err)
	}
}
