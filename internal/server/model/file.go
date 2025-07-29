package model

import (
	"time"

	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/service"
)

type File struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type FileRequest struct {
	UserID      int64
	Name        string
	AccessToken string
}

func GetListFile(binary []File) []*grpc.Binary {
	items := make([]*grpc.Binary, len(binary))
	for i := range binary {
		created := service.ConvertTimeToTimestamp(binary[i].CreatedAt)
		items[i] = &grpc.Binary{Id: binary[i].ID, Name: binary[i].Name, CreatedAt: created}
	}
	return items
}
