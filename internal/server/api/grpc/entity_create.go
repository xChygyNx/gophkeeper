package grpchandler

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpc "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/errors"
)

// EntityCreate - check the validity of the token and save record (text, bank card or login password)
func (h *Handler) EntityCreate(ctx context.Context, req *grpc.CreateEntityRequest) (*grpc.CreateEntityResponse, error) {
	h.logger.Info("entity create")

	logError, statusErr := validateAccessToken(req.AccessToken.Token, h.token)
	if statusErr != nil {
		h.logger.Error(logError)
		return &grpc.CreateEntityResponse{}, statusErr
	}

	var metadata model.MetadataEntity
	err := json.Unmarshal([]byte(req.Metadata), &metadata)
	if err != nil {
		finalError := fmt.Errorf("error in Unmarshar request metadata: %w", err)
		h.logger.Error(finalError)
		return &grpc.CreateEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}

	EntityData := &model.CreateEntityRequest{}
	EntityData.UserID = req.AccessToken.UserId
	EntityData.Data = req.Data
	EntityData.Metadata = metadata
	if metadata.Name == "" {
		err := errors.ErrNoMetadataSet
		h.logger.Error(err)
		return &grpc.CreateEntityResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
	}

	exists, err := h.entity.Exists(EntityData)
	if err != nil {
		finalError := fmt.Errorf("error in check of entity exists: %w", err)
		h.logger.Error(finalError)
		return &grpc.CreateEntityResponse{}, status.Errorf(codes.Internal, finalError.Error())
	}
	if exists {
		err = errors.ErrNameAlreadyExists
		h.logger.Error(err)
		return &grpc.CreateEntityResponse{}, status.Errorf(codes.AlreadyExists, err.Error())
	}

	CreatedEntityID, err := h.entity.Create(EntityData)
	if err != nil {
		h.logger.Error(err)
		return &grpc.CreateEntityResponse{}, status.Errorf(
			codes.Internal, err.Error(),
		)
	}

	h.logger.Debug("Created entity ID: %d", CreatedEntityID)
	return &grpc.CreateEntityResponse{Id: CreatedEntityID}, nil
}
