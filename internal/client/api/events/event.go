package events

import (
	"context"
	grpc "github.com/xChygyNx/gophkeeper/internal/proto"

	"github.com/sirupsen/logrus"
	"github.com/xChygyNx/gophkeeper/internal/client/config"
)

type Event struct {
	grpc    grpc.GophkeeperClient
	config  *config.ConfigClient
	logger  *logrus.Logger
	context context.Context
	grpc.UnimplementedGophkeeperServer
}

// NewEvent - creates a new grpc client instance
func NewEvent(ctx context.Context, config *config.ConfigClient, log *logrus.Logger, client grpc.GophkeeperClient) *Event {
	return &Event{context: ctx, config: config, logger: log, grpc: client}
}

// GetConfig - returns config in string format
func (s Event) GetConfig() *config.ConfigClient {
	return s.config
}
