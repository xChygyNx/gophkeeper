package repositories

import (
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	"github.com/xChygyNx/gophkeeper/internal/server/storage"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/entity"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/file"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/user"
)

type Repositories struct {
	UserRepo   *user.User
	BinaryRepo *file.File
	AppStorage *storage.Storage
	EntityRepo *entity.Entity
	TokenRepo  *token.Token
}

func InitRepositories(db *database.DB, appStoragePath string) *Repositories {
	appStorage := storage.New(appStoragePath)
	return &Repositories{
		UserRepo:   user.New(db),
		BinaryRepo: file.New(db),
		AppStorage: &appStorage,
		EntityRepo: entity.New(db),
		TokenRepo:  token.New(db),
	}

}
