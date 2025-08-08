package model

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/timestamp"

	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

type User struct {
	ID        int64
	Username  string
	Password  string
	CreatedAt timestamp.Timestamp
	UpdatedAt timestamp.Timestamp
	DeletedAt timestamp.Timestamp
}

type GetAllUsers struct {
	ID        int64
	Username  string
	Password  string
	DeletedAt sql.NullTime
}

type UserRequest struct {
	Username string
	Password string
}

func GetUserData(data *User) grpc.User {
	return grpc.User{
		UserId:    data.ID,
		Username:  data.Username,
		CreatedAt: &data.CreatedAt,
		UpdatedAt: &data.UpdatedAt,
		DeletedAt: &data.DeletedAt,
	}
}
