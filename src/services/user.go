package services

import (
	"context"
	"database/sql"
	"fmt"
	"grpc-go/src/config"
	userPb "grpc-go/src/pb/user"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserService struct {
	userPb.UnimplementedUserServiceServer
	DB *config.DatabaseConfig
}

func (userService *UserService) ListUsers(context.Context, *userPb.Empty) (result *userPb.ListUsersResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	if err != nil {
		fmt.Println("begin error", err.Error())
		result = nil
		return result, nil
	}
	var rows *sql.Rows
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
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	var rows *sql.Rows
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
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	var rows *sql.Rows
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
func (userService *UserService) CreateUser(ctx context.Context, toCreateUser *userPb.CreateUserRequest) (result *userPb.CreateUserResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	userId := uuid.New()
	time := time.Now()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(toCreateUser.Password), bcrypt.DefaultCost)
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("begin error", err.Error())
		result = nil
		return result, rollback
	}
	_, err = begin.Query(
		`INSERT INTO "users" (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);`,
		userId,
		toCreateUser.Name,
		toCreateUser.Email,
		hashedPassword,
		time,
		time,
	)
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("query error", err.Error())
		result = nil
		return result, rollback
	}
	commit := begin.Commit()
	response := &userPb.CreateUserResponse{
		Code:    int64(codes.OK),
		Message: "CreateUser Succes",
		Data: &userPb.User{
			Id:        userId.String(),
			Name:      toCreateUser.Name,
			Email:     toCreateUser.Email,
			Password:  string(hashedPassword),
			CreatedAt: timestamppb.New(time),
			UpdatedAt: timestamppb.New(time),
		},
	}
	return response, commit
}
func (userService *UserService) UpdateUser(ctx context.Context, toUpdateUser *userPb.UpdateUserRequest, userId *userPb.ById) (result *userPb.UpdateUserResponse, err error) {
	begin, err := userService.DB.GrpcDB.Connection.Begin()
	if err != nil {
		result = &userPb.UpdateUserResponse{
			Message: "begin failed, " + err.Error(),
		}
		rollback := begin.Rollback()
		return result, rollback
	}
	time := time.Now()
	var hashedPassword []byte
	if toUpdateUser.Password != nil {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*toUpdateUser.Password), bcrypt.DefaultCost)
		if err != nil {
			result = &userPb.UpdateUserResponse{
				Message: "UpdateUser failed, hashing password fail, " + err.Error(),
			}
			rollback := begin.Rollback()
			return result, rollback
		}
	}
	rows, err := begin.Query(
		`SELECT id, name, email, password, balance, created_at, updated_at, deleted_at FROM "users" WHERE id=$1 LIMIT 1;`,
		userId,
	)

	if err != nil {
		result = &userPb.UpdateUserResponse{
			Message: "UpdateUser failed, query GetUserById fail, " + err.Error(),
		}
		rollback := begin.Rollback()
		return result, rollback
	}
	defer rows.Close()
	_, err = begin.Query(
		`INSERT INTO "users" (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);`,
		userId,
		toUpdateUser.Name,
		toUpdateUser.Email,
		hashedPassword,
		time,
		time,
	)
	if err != nil {
		rollback := begin.Rollback()
		fmt.Println("query error", err.Error())
		result = nil
		return result, rollback
	}
	commit := begin.Commit()
	response := &userPb.UpdateUserResponse{
		Code:    int64(codes.OK),
		Message: "CreateUser Succes",
		Data: &userPb.User{
			Id:        userId.String(),
			Name:      toUpdateUser.Name,
			Email:     toUpdateUser.Email,
			Password:  string(hashedPassword),
			CreatedAt: timestamppb.New(time),
			UpdatedAt: timestamppb.New(time),
		},
	}
	return response, commit
}
