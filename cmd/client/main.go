package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/client/api/events"
	"github.com/xChygyNx/gophkeeper/internal/client/config"
	"github.com/xChygyNx/gophkeeper/internal/client/gui"
	gophkeeper "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// InterceptorLogger adapts logrus logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l logrus.FieldLogger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make(map[string]any, len(fields)/2)
		i := logging.Fields(fields).Iterator()
		if i.Next() {
			k, v := i.At()
			f[k] = v
		}
		l := l.WithFields(f)

		switch lvl {
		case logging.LevelDebug:
			l.Debug(msg)
		case logging.LevelInfo:
			l.Info(msg)
		case logging.LevelWarn:
			l.Warn(msg)
		case logging.LevelError:
			l.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

// @Title Password Manager gophkeeper
// @Description GophKeeper is a client-server system that allows the user to safely and securely store logins, passwords, binary data and other private information.
// @Version 1.0

// @Contact.email pavel@utkin-pro.ru

func main() {
	logger := logrus.New()
	ctx := context.Background()
	clientConfig := config.NewConfig()

	conn, err := grpc.Dial(
		clientConfig.GRPC,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(InterceptorLogger(logger)),
		),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(InterceptorLogger(logger)),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	gophkeeperClient := gophkeeper.NewGophkeeperClient(conn)
	client := events.NewEvent(ctx, clientConfig, logger, gophkeeperClient)
	_, err = client.Ping()
	if err != nil {
		log.Fatal(err)
	}
	application := app.New()
	gui.InitGUI(logger, application, client)
}
