package main

import (
	"log"
	"net"

	"grpc-go/src/config"
	userPb "grpc-go/src/pb/user"
	"grpc-go/src/services"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := ":50051"
	netListen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	log.Printf("server started at %v", port)

	grpcServer := grpc.NewServer()
	envConfig := config.NewEnvConfig()
	db := config.NewGrpcDBConfig(envConfig)

	userService := services.UserService{DB: db}

	userPb.RegisterUserServiceServer(grpcServer, &userService)

	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatalf("failed to serve %v", err.Error())
	}
}
