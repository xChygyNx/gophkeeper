package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/server/api"
	grpchandler "github.com/xChygyNx/gophkeeper/internal/server/api/grpc"
	resthandler "github.com/xChygyNx/gophkeeper/internal/server/api/rest"
	"github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	"github.com/xChygyNx/gophkeeper/internal/server/storage"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/entity"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/file"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/user"
)

// @Title Password Manager github.com/xChygyNx/gophkeeper
// @Description GophKeeper is a client-server system that allows the user to safely and securely store logins, passwords, binary data and other private information.
// @Version 1.0

// @Contact.email pavel@utkin-pro.ru

func main() {
	logger := logrus.New()
	serverConfig, err := config.NewConfig(logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetLevel(serverConfig.DebugLevel)

	db, err := database.New(serverConfig, logger)
	if err != nil {
		logger.Fatal(err)
	} else {
		defer func() {
			err = db.Close()
			if err != nil {
				logger.Fatalf("db close failed: %v", err)
			}
		}()
		err = db.CreateTablesMigration("file://migrations")
		if err != nil {
			logger.Fatalf("Migration failed: %v", err)
		}
	}

	userRepository := user.New(db)
	binaryRepository := file.New(db)
	appStorage := storage.New("/tmp")
	entityRepository := entity.New(db)
	tokenRepository := token.New(db)

	handlerRest := resthandler.NewHandler(db, serverConfig, userRepository, tokenRepository, logger)
	routerService := resthandler.Route(handlerRest)
	rs := chi.NewRouter()
	rs.Mount("/", routerService)

	handlerGrpc := grpchandler.NewHandler(db, serverConfig, userRepository, binaryRepository,
		&appStorage, entityRepository, tokenRepository, logger)

	ctx, cnl := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cnl()

	go api.StartGRPCService(handlerGrpc, serverConfig, logger)
	go api.StartRESTService(rs, serverConfig, logger)

	<-ctx.Done()
	logger.Info("server shutdown on signal with:", ctx.Err())
}
