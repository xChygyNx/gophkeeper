package grpchandler

import (
	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/entity"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/file"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/user"
)

type Handler struct {
	database *database.DB
	config   *config.Config
	user     *user.User
	file     *file.File
	storage  *storage.Storage
	entity   *entity.Entity
	token    *token.Token
	logger   *logrus.Logger
	grpc.UnimplementedGophkeeperServer
}

// NewHandler - creates a new grpc server instance
func NewHandler(db *database.DB, config *config.Config, userRepository *user.User,
	binaryRepository *file.File, storage *storage.Storage, entityRepository *entity.Entity, tokenRepository *token.Token, log *logrus.Logger) *Handler {
	return &Handler{database: db, config: config, user: userRepository, file: binaryRepository, storage: storage,
		entity: entityRepository, token: tokenRepository, logger: log}
}
