package main

import (
	pb "biz-server/biz"
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

type postgres struct {
	db *pgxpool.Pool
}

var (
	port = flag.Int("port", 5062, "The server port")
	//db    *gorm.DB
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

type getUserServer struct {
	pb.UnimplementedGet_UsersServer
}

func (pg *postgres) GetUsers(ctx context.Context,
	in *pb.Get_Users_Req) (*pb.Get_Users_Resp, error) {
	fmt.Println("Get user request received")
	var usr []*pb.USERS
	if &in.UserId != nil {
		queryString := "SELECT * FROM USERS WHERE user_id=" + string(in.GetUserId()) + ";"

		rows, err := pg.db.Query(context.Background(), queryString)
		if err != nil {
			fmt.Printf("rows.Scan error: %s", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			rowValues, _ := rows.Values()
			usr = append(usr, &pb.USERS{
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
			fmt.Printf("rows.Scan error: %s", err)
			return nil, err
		}
		defer rows.Close()
		itterate := 0
		for rows.Next() && itterate < 100 {
			itterate++
			rowValues, _ := rows.Values()
			usr = append(usr, &pb.USERS{
				Id:        rowValues[0].(int64),
				Name:      rowValues[1].(string),
				Family:    rowValues[2].(string),
				Age:       rowValues[3].(int64),
				Sex:       rowValues[4].(string),
				CreatedAt: rowValues[5].(string),
			})
		}
	}
	return &pb.Get_Users_Resp{Users: usr, MessageId: in.GetMessageId() + 1}, nil
}

func newGetUserServer() *getUserServer {
	return &getUserServer{}
}

func (pg *postgres) GetUserWSqlInj(ctx context.Context,
	in *pb.Get_User_Sql_Inj_Req) (*pb.Get_Users_Resp, error) {
	fmt.Println("Get user with sql injection request received")
	queryString := "SELECT * FROM USERS WHERE user_id=" + in.GetUserId() + ";"
	fmt.Println("query string: " + queryString)

	_, err := pg.db.Query(context.Background(), queryString)
	if err != nil {
		fmt.Printf("rows.Scan error: %s", err)
		return nil, err
	}
	return &pb.Get_Users_Resp{Users: nil, MessageId: in.GetMessageId() + 1}, nil
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

	//db, err := sql.Open("postgres", psqlconn)
	db, err := pgxpool.Connect(context.Background(), psqlconn)

	checkError(err)

	fmt.Println("connected to database!")
	defer db.Close()
	queryCreateTable := "Create table USERS(\n\tuser_id int,\n\tname varchar,\n\tfamily varchar,\n\tage int,\n\tsex varchar(4),\n\tcreatedAt varchar\n);"
	_, err = db.Exec(context.Background(), queryCreateTable)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create USERS table: %v\n", err)
		log.Fatal(err)
	}
	for i := range users {
		queryInsertUsers := `INSERT INTO USERS (user_id, name, family, age, 
                   createdAt, sex) VALUES ($1, $2,$3,$4,$5,$6);`
		_, err := db.Exec(context.Background(), queryInsertUsers, users[i].GetId(), users[i].GetName(),
			users[i].GetFamily(), users[i].GetAge(), users[i].GetCreatedAt(), users[i].GetSex())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", err)
			log.Fatal(err)
		}
	}

	fmt.Println("insertion completed!")

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
