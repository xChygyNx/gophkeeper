package config

import (
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	resthandler "github.com/xChygyNx/gophkeeper/internal/server/api/rest"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
)

func InitHTTPRouter(db *database.DB, serverConfig *Config, repositories *Repositories,
	logger *logrus.Logger) *chi.Mux {
	handlerRest := resthandler.NewHandler(db, serverConfig, repositories.UserRepo, repositories.TokenRepo, logger)
	routerService := resthandler.Route(handlerRest)
	rs := chi.NewRouter()
	rs.Mount("/", routerService)

	return rs
}
