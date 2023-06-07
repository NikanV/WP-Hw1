package main

import (
	pb "biz-server/biz"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

type postgres struct {
	db *pgxpool.Pool
}

type bizServiceServer struct {
	pb.UnimplementedBizServiceServer
}

func (c *bizServiceServer) GetUsers(ctx context.Context, in *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users_response := []*pb.USER{}
	user_id := in.GetUserId()
	message_id := in.GetMessageId()
	if user_id >= 0 {
		queryString := "SELECT * FROM USERS WHERE user_id=" + strconv.Itoa(int(user_id)) + ";"

		rows, err := pg.db.Query(context.Background(), queryString)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			rowValues, _ := rows.Values()
			users_response = append(users_response, &pb.USER{
				Id:        rowValues[0].(int64),
				Name:      rowValues[1].(string),
				Family:    rowValues[2].(string),
				Age:       rowValues[3].(int64),
				Sex:       rowValues[4].(string),
				CreatedAt: rowValues[5].(string),
			})
		}
	} else {
		queryString := "SELECT * FROM USERS;"
		rows, err := pg.db.Query(context.Background(), queryString)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for iterate := 0; rows.Next() && iterate < 100; iterate++ {
			iterate++
			rowValues, _ := rows.Values()
			users_response = append(users_response, &pb.USER{
				Id:        rowValues[0].(int64),
				Name:      rowValues[1].(string),
				Family:    rowValues[2].(string),
				Age:       rowValues[3].(int64),
				Sex:       rowValues[4].(string),
				CreatedAt: rowValues[5].(string),
			})
		}
	}

	response := pb.GetUsersResponse{
		Users:     users_response,
		MessageId: message_id + 1,
	}

	return &response, nil
}

func (c *bizServiceServer) GetUsersWithSQL(ctx context.Context, in *pb.GetUsersWithSQLRequest) (*pb.GetUsersResponse, error) {
	user_id := in.GetUserId()
	message_id := in.GetMessageId()
	queryString := "SELECT * FROM USERS WHERE user_id=" + user_id + ";"

	rows, err := pg.db.Query(context.Background(), queryString)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rowValues, _ := rows.Values()
	users_response := []*pb.USER{}
	users_response[0] = &pb.USER{
		Id:        rowValues[0].(int64),
		Name:      rowValues[1].(string),
		Family:    rowValues[2].(string),
		Age:       rowValues[3].(int64),
		Sex:       rowValues[4].(string),
		CreatedAt: rowValues[5].(string),
	}

	response := pb.GetUsersResponse{
		Users:     users_response,
		MessageId: message_id + 1,
	}

	return &response, nil
}

func newBizServiceServer() *bizServiceServer {
	return &bizServiceServer{}
}

var (
	port  = flag.Int("port", 5062, "The server port")
	users = []*pb.USER{
		{Name: "nikan", Family: "vsi", Id: 5303, Age: 19, Sex: "male", CreatedAt: time.Now().UTC().String()},
		{Name: "nima", Family: "enigma", Id: 5263, Age: 19, Sex: "male", CreatedAt: time.Now().UTC().String()},
		{Name: "alireza", Family: "kaz", Id: 9649, Age: 19, Sex: "male", CreatedAt: time.Now().UTC().String()},
	}
	pg = postgres{}
)

const (
	db_host = "localhost"
	db_port = 5432
	db_user = "postgres"
	db_pass = "postgres"
	db_name = "postgres"
)

func main() {
	flag.Parse()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_pass, db_name)

	opened_db, err := pgxpool.Connect(context.Background(), psqlconn)
	if err != nil {
		log.Fatalf("Failed to connect to database! %v", err)
	} else {
		pg.db = opened_db
		fmt.Println("Connected to database!")
	}
	defer pg.db.Close()

	queryCreateTable := "Create table USERS(\n\tuser_id int,\n\tname varchar,\n\tfamily varchar,\n\tage int,\n\tsex varchar(4),\n\tcreatedAt varchar\n);"
	if _, err = pg.db.Exec(context.Background(), queryCreateTable); err != nil {
		log.Fatalf("Unable to create USERS table! %v", err)
	}
	for i := range users {
		queryInsertUsers := `INSERT INTO USERS (user_id, name, family, age, 
                   createdAt, sex) VALUES ($1, $2,$3,$4,$5,$6);`
		if _, err := pg.db.Exec(context.Background(), queryInsertUsers, users[i].GetId(), users[i].GetName(), users[i].GetFamily(), users[i].GetAge(), users[i].GetCreatedAt(), users[i].GetSex()); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", err)
			log.Fatal(err)
		}
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen! %v", err)
	} else {
		fmt.Println("Listening on port: ", *port)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBizServiceServer(grpcServer, newBizServiceServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve! %v", err)
	}
}
