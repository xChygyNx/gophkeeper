package events

import (
	"fmt"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
)

// UserExist - check if user exist in db
func (s Event) UserExist(username string) (bool, error) {
	s.logger.Info("user exist check")

	user, err := s.grpc.UserExist(s.context, &grpc.UserExistRequest{Username: username})
	if err != nil {
		myErr := fmt.Errorf("error in check of user exists: %w", err)
		s.logger.Error(myErr)
		return user.Exist, myErr
	}

	return user.Exist, nil
}
