package services

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-go/src/config"
	userPb "grpc-go/src/pb/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	userPb.UnimplementedUserServiceServer
	DB *config.DatabaseConfig
}

func (userService *UserService) ListUsers(context.Context, *userPb.Empty) (result *userPb.ListUsersResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	var rows *sql.Rows
	if err != nil {
		fmt.Println("begin error", err.Error())
		result = nil
		return result, nil
	}
	var ListUsers []*userPb.User

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
		ListUser := &userPb.User{}
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
		ListUser.CreatedAt = timestamppb.New(createdAt)
		ListUser.UpdatedAt = timestamppb.New(createdAt)

		ListUsers = append(ListUsers, ListUser)
	}
	commit := begin.Commit()
	response := &userPb.ListUsersResponse{
		Code:    int64(codes.OK),
		Message: "ListUser Succeed",
		Data:    ListUsers,
	}
	return response, commit
}
func (userService *UserService) DeleteUser(_ context.Context, id *userPb.ById) (result *userPb.DeleteUserResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	var rows *sql.Rows
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	rows, err = begin.Query(
		`DELETE FROM "users" WHERE id=$1 RETURNING id, name,  email, password, created_at, updated_at`,
		id.Id,
	)
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("query error", err.Error())
		result = nil
		return result, rollback
	}
	defer rows.Close()
	var UserData []*userPb.User
	for rows.Next() {
		user := &userPb.User{}
		createdAt := user.CreatedAt.AsTime()
		updatedAt := user.UpdatedAt.AsTime()
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			rollback := begin.Rollback()
			fmt.Println("scan error", err.Error())
			result = nil
			return result, rollback
		}
		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(createdAt)
		UserData = append(UserData, user)
	}
	commit := begin.Commit()
	response := &userPb.DeleteUserResponse{
		Code:    int64(codes.OK),
		Message: "DeleteUser Succeed",
		Data:    UserData[0],
	}
	return response, commit
}

func (userService *UserService) DetailUser(_ context.Context, id *userPb.ById) (result *userPb.DetailUserResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	var rows *sql.Rows
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	rows, err = begin.Query(
		`SELECT id, name, email, password, created_at, updated_at FROM "users" WHERE id=$1 LIMIT 1;`,
		id.Id,
	)
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("query error", err.Error())
		result = nil
		return result, rollback
	}
	defer rows.Close()
	var UserData []*userPb.User
	for rows.Next() {
		user := &userPb.User{}
		createdAt := user.CreatedAt.AsTime()
		updatedAt := user.UpdatedAt.AsTime()
		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			rollback := begin.Rollback()
			fmt.Println("scan error", err.Error())
			result = nil
			return result, rollback
		}
		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(createdAt)
		UserData = append(UserData, user)
	}
	commit := begin.Commit()
	response := &userPb.DetailUserResponse{
		Code:    int64(codes.OK),
		Message: "ListUser Succeed",
		Data:    UserData[0],
	}
	return response, commit
}
