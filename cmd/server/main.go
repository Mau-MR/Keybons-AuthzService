package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Mau-MR/authzService/pb"
	"github.com/Mau-MR/authzService/service"
	"google.golang.org/grpc"
)

const (
	tokenDuration = 5 * time.Minute
	secretKey     = "secretjaja"
)

func main() {
	dbClient := service.NewDataBase("localhost", 27017)
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port: %d", *port)
	//To manage the jwt
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(dbClient, jwtManager)
	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
