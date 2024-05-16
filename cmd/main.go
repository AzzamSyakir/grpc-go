package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	port := ":50051"
	netListen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	log.Printf("server started at %v", port)
	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatalf("failed to serve %v", err.Error())
	}

}
