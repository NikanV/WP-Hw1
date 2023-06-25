package main

import (
	"context"
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

func (c *authServer) AuthCheck(ctx context.Context, in *pb.ACRequest) (*pb.ACResponse, error) {
	message_id := in.GetMessageId()
	auth_key := in.GetAuthKey()

	client := initRedisClient()
	defer client.Close()

	auth_check, _ := client.Exists(context.Background(), strconv.Itoa(int(auth_key))).Result()

	response := pb.ACResponse{
		MessageId: message_id + 1,
		AuthCheck: auth_check > 0,
	}

	return &response, nil
}

func (c *authServer) RequestPQ(ctx context.Context, in *pb.PQRequest) (*pb.PQResponse, error) {
	message_id := in.GetMessageId()
	nonce := in.GetNonce()

	prime := tools.Random_Prime()
	p, g := int64(prime), int64(tools.FindPrimitive(prime))
	server_nonce := tools.RandomString(20)

	client := initRedisClient()
	defer client.Close()
	hash := tools.Sha1_gen(nonce + server_nonce)
	err := client.HSet(context.Background(), hash, "p", p).Err()
	if err != nil {
		return nil, err
	}
	err = client.HSet(context.Background(), hash, "g", g).Err()
	if err != nil {
		return nil, err
	}
	client.Expire(context.Background(), hash, 20*time.Minute)

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

	private_key := tools.RandomNumber(10)
	hash := tools.Sha1_gen(nonce + server_nonce)
	client := initRedisClient()
	defer client.Close()

	redis_p, err := client.HGet(context.Background(), hash, "p").Result()
	if err == redis.Nil {
		return nil, err
	}
	redis_g, err := client.HGet(context.Background(), hash, "g").Result()
	if err == redis.Nil {
		return nil, err
	}
	p, _ := strconv.ParseInt(redis_p, 10, 64)
	g, _ := strconv.ParseInt(redis_g, 10, 64)
	auth_key := int64(p) % int64(math.Pow(float64(a), float64(private_key)))

	client.Del(context.Background(), hash)
	err = client.Set(context.Background(), strconv.Itoa(int(auth_key)), 0, 0).Err()
	if err != nil {
		return nil, err
	}

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
		Addr:     "host.docker.internal:6379",
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

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
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
