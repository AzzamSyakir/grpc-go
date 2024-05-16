package services

import (
	"context"
	"database/sql"
	"grpc-go/src/config"
	"grpc-go/src/entity"
	userPb "grpc-go/src/pb/user"
	"log"
)

type UserService struct {
	userPb.UnimplementedUserServiceServer
	DB *config.DatabaseConfig
}

func (u *UserService) ListUsers(context.Context, *userPb.Empty) (result *userPb.User, err error) {
	begin, err := u.DB.GrpcDB.Connection.Begin()
	var rows *sql.Rows
	if err != nil {
		rollback := begin.Rollback()
		log.Printf("begin error %v", err.Error())
		result = nil
		return result, rollback
	}
	rows, err = begin.Query(
		`SELECT id, name, email, password, created_at, updated_at, deleted_at FROM "users" `,
	)
	defer rows.Close()
	var ListUsers []*entity.User
	for rows.Next() {
		ListUser := &entity.User{}
		scanErr := rows.Scan(
			&ListUser.Id,
			&ListUser.Name,
			&ListUser.Email,
			&ListUser.Password,
			&ListUser.CreatedAt,
			&ListUser.UpdatedAt,
			&ListUser.DeletedAt,
		)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		ListUsers = append(ListUsers, ListUser)
	}
	if err != nil {
		rollback := begin.Rollback()
		log.Printf("query error %v", err.Error())
		result = nil
		return result, rollback
	}
	return ListUsers, nil
}
