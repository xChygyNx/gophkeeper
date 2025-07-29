package resthandler

import (
	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/user"
)

type Handler struct {
	database *database.DB
	config   *config.Config
	user     *user.User
	token    *token.Token
	log      *logrus.Logger
}

// NewHandler - creates a new server instance
func NewHandler(db *database.DB, config *config.Config, userRepository *user.User, tokenRepository *token.Token,
	log *logrus.Logger) *Handler {
	return &Handler{database: db, config: config, user: userRepository, token: tokenRepository, log: log}
}
