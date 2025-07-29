package events

import (
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// Ping - ping
func (s Event) Ping() (string, error) {
	s.logger.Info("ping")

	msg, err := s.grpc.Ping(s.context, &grpc.PingRequest{})
	if err != nil {
		myErr := fmt.Errorf("error in ping server: %w", err)
		s.logger.Error(myErr)
		return "", myErr
	}

	return msg.Message, nil
}
