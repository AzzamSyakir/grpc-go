package main

import (
	"fmt"
	pb "grpc-go/pb"
	"log"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func main() {
	id := uuid.New()
	user := &pb.User{
		Id:       id.String(),
		Name:     "asa",
		Email:    "asa@gmail.com",
		Password: "password123",
	}
	data, err := proto.Marshal(user)
	if err != nil {
		log.Fatal("marshaling error, ", err)
	}
	fmt.Println("success marshal", data)
	testUser := &pb.User{}
	if err = proto.Unmarshal(data, testUser); err != nil {
		log.Fatal("unmarshaling error, ", err)
	}
	fmt.Println("succes unmarshal", testUser)
}
