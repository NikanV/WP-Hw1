package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	pb "WP-Hw1/proto"
	tools "WP-Hw1/tools"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type authServer struct {
	pb.UnimplementedAuthenticatorServer
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
	values, _ := json.Marshal(map[string]int64{
		"p": p,
		"g": g,
	})
	hash := tools.Sha1_gen(nonce + server_nonce)

	err := client.HSet(context.Background(), hash, values, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

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
	private_key := tools.RandomNumber(50)
	hash := tools.Sha1_gen(nonce + server_nonce)
	client := initRedisClient()
	defer client.Close()
	data := client.HGetAll(context.Background(), hash).Val()

	p, _ := strconv.ParseInt(data["p"], 10, 64)
	g, _ := strconv.ParseInt(data["g"], 10, 64)
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
		MessageId:   message_id,
		B:           int64(math.Pow(float64(p), float64(private_key))) % int64(g),
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
