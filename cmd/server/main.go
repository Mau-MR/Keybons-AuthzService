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
	"google.golang.org/grpc/reflection"
)

//The constants
const (
	tokenDuration = 5 * time.Minute
	secretKey     = "secretjaja"
)

//TODO: Check for the authentication of every service on another specific package
//This stands for all the routes that this service redirects
func accesibleRoles() map[string][]string {
	const clientServicePath = "/pb.AuthService/"
	return map[string][]string{
		clientServicePath + "AccountExistance": {"admin", "user", "owner"},
		clientServicePath + "CreateClient":     {"admin"},
		clientServicePath + "UploadImage":      {"user", "admin"},
	}
}

func main() {
	//The port configuration
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port: %d", *port)

	//The service configuration part
	dbClient := service.NewDataBase("localhost", 27017)
	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(dbClient, jwtManager)
	//interceptor
	interceptor := service.NewAuthInterceptor(jwtManager, accesibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	//The register of the services on the server
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	//This part is for making some test during development
	reflection.Register(grpcServer)

	//Making the connections
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}
