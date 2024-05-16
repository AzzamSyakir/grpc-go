package services

import (
	"context"
	userPb "grpc-go/src/pb/user"

	"github.com/google/uuid"
)

type UserService struct {
	userPb.UnimplementedUserServiceServer
}

func (u *UserService) ListUsers(context.Context, *userPb.Empty) (*userPb.User, error) {
	uuid := uuid.New()
	userData := &userPb.User{
		Id:       uuid.String(),
		Name:     "name" + uuid.String(),
		Email:    "email" + uuid.String(),
		Password: "password" + uuid.String(),
	}
	return userData, nil
}
