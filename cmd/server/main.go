package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/server/api"
	grpchandler "github.com/xChygyNx/gophkeeper/internal/server/api/grpc"
	"github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
)

// @Title Password Manager github.com/xChygyNx/gophkeeper
// @Description GophKeeper is a client-server system that allows the user to safely and securely store logins, passwords, binary data and other private information.
// @Version 1.0

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

	repositories := config.InitRepositories(db, "/tmp")

	rs := config.InitHTTPRouter(db, serverConfig, repositories, logger)

	handlerGrpc := grpchandler.NewHandler(db, serverConfig, repositories.UserRepo, repositories.BinaryRepo,
		repositories.AppStorage, repositories.EntityRepo, repositories.TokenRepo, logger)

	ctx, cnl := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cnl()

	go api.StartGRPCService(handlerGrpc, serverConfig, logger)
	go api.StartRESTService(rs, serverConfig, logger)

	<-ctx.Done()
	logger.Info("server shutdown on signal with:", ctx.Err())
}
