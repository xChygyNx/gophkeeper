package resthandler

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories"

	"github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
)

func InitHTTPRouter(db *database.DB, serverConfig *config.Config, repositories *repositories.Repositories,
	logger *logrus.Logger) *chi.Mux {
	handlerRest := NewHandler(db, serverConfig, repositories.UserRepo, repositories.TokenRepo, logger)
	routerService := Route(handlerRest)
	rs := chi.NewRouter()
	rs.Mount("/", routerService)

	return rs
}
