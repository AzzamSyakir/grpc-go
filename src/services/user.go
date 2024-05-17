package services

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-go/src/config"
	userPb "grpc-go/src/pb/user"

	"google.golang.org/grpc/codes"
)

type UserService struct {
	userPb.UnimplementedUserServiceServer
	DB *config.DatabaseConfig
}

func (userService *UserService) ListUsers(context.Context, *userPb.Empty) (result *userPb.ListUserResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	var rows *sql.Rows
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	var ListUsers []*userPb.UserResponse

	rows, err = begin.Query(
		`SELECT id, name, email, password,  created_at, updated_at FROM users`,
	)
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("query error", err.Error())
		result = nil
		return result, rollback
	}
	defer rows.Close()
	for rows.Next() {
		ListUser := &userPb.UserResponse{}
		createdAt := ListUser.CreatedAt.AsTime()
		updatedAt := ListUser.UpdatedAt.AsTime()
		err = rows.Scan(
			&ListUser.Id,
			&ListUser.Name,
			&ListUser.Email,
			&ListUser.Password,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			rollback := begin.Rollback()
			fmt.Println("scan error", err.Error())
			result = nil
			return result, rollback
		}
		ListUsers = append(ListUsers, ListUser)
	}
	commit := begin.Commit()
	response := &userPb.ListUserResponse{
		Code:    int64(codes.OK),
		Message: "ListUser Succeed",
		Data:    ListUsers,
	}
	return response, commit
}
