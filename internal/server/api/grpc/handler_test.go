package grpchandler

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/client/service/randomizer"
	serverConfig "github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	"github.com/xChygyNx/gophkeeper/internal/server/model"
	grpcKeeper "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/entity"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/file"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/user"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/vars"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestHandlers(t *testing.T) {

	// -- SETUP --
	// initiate postgres container
	container, err := postgres.RunContainer(context.Background(),
		testcontainers.WithImage("docker.io/postgres:latest"),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("Test containers failed: %v", err)
	}

	err = container.Start(context.Background())
	if err != nil {
		t.Fatalf("Test containers failed: %v", err)
	}

	stopTime := time.Second
	defer func() {
		err = container.Stop(context.Background(), &stopTime)
		if err != nil {
			t.Fatalf("container stop failed: %v", err)
		}
	}()

	databaseURI, err := container.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		t.Fatalf("container connection failed: %v", err)
	}
	logger := logrus.New()
	db, err := database.New(&serverConfig.Config{DSN: databaseURI}, logger)
	if err != nil {
		t.Fatalf("Db init failed: %v", err)
	}

	err = db.CreateTablesMigration("file://../../../../migrations")
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	serverCnfg := &serverConfig.Config{
		AddressGRPC:         "localhost:8080",
		DSN:                 databaseURI,
		AccessTokenLifetime: 300 * time.Second,
		DebugLevel:          logrus.DebugLevel,
		FileFolder:          "../../../../data/test_keeper",
	}

	// repositories
	userRepository := user.New(db)
	fileRepository := file.New(db)
	storage := storage.New("/tmp")
	entityRepository := entity.New(db)
	tokenRepository := token.New(db)

	// setup server service
	handlerGrpc := *NewHandler(db, serverCnfg, userRepository, fileRepository, &storage,
		entityRepository, tokenRepository, logger)

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	t.Run("registration gophkeeper-server", func(t *testing.T) {
		grpcKeeper.RegisterGophkeeperServer(s, &handlerGrpc)

		go func() {
			if err := s.Serve(lis); err != nil {
				t.Errorf("Server exited with error: %v", err)
				return
			}
		}()
	})

	// -- TEST DATA --
	var authenticatedUser *grpcKeeper.AuthenticationResponse
	var blockedUser *grpcKeeper.AuthenticationResponse
	data := randomizer.RandStringRunes(10)
	dataUpdate := randomizer.RandStringRunes(10)
	username := randomizer.RandStringRunes(10)
	password, _ := encryption.HashPassword("Password-00")

	blockedusername := randomizer.RandStringRunes(10)
	blockedpassword, _ := encryption.HashPassword("Password-00")

	name := randomizer.RandStringRunes(10)
	description := randomizer.RandStringRunes(10)
	metadata := model.MetadataEntity{Name: name, Description: description, Type: vars.Text.ToString()}
	jsonMetadata, _ := json.Marshal(metadata)

	// -- TESTS --
	t.Run("ping db", func(t *testing.T) {
		_, err = handlerGrpc.Ping(context.Background(), &grpcKeeper.PingRequest{})
		assert.NoError(t, err, "failed ping db")
	})

	t.Run("registration", func(t *testing.T) {
		_, err = handlerGrpc.Registration(context.Background(), &grpcKeeper.RegistrationRequest{Username: username, Password: password})
		assert.NoError(t, err, "registration failed")
	})

	t.Run("registration", func(t *testing.T) {
		_, err = handlerGrpc.Registration(context.Background(),
			&grpcKeeper.RegistrationRequest{Username: blockedusername, Password: blockedpassword})
		assert.NoError(t, err, "registration failed")
	})

	t.Run("registration", func(t *testing.T) {
		_, err = handlerGrpc.Registration(context.Background(), &grpcKeeper.RegistrationRequest{Username: username, Password: password})
		assert.Error(t, err, "registration failed")
	})

	t.Run("user exist", func(t *testing.T) {
		_, err = handlerGrpc.UserExist(context.Background(), &grpcKeeper.UserExistRequest{Username: username})
		assert.NoError(t, err, "user exist failed")
	})

	t.Run("authentication", func(t *testing.T) {
		authenticatedUser, err = handlerGrpc.Authentication(context.Background(), &grpcKeeper.AuthenticationRequest{Username: username, Password: password})
		assert.NoError(t, err, "authentication failed")
	})

	t.Run("authentication and block user", func(t *testing.T) {
		blockedUser, err = handlerGrpc.Authentication(context.Background(),
			&grpcKeeper.AuthenticationRequest{Username: blockedusername, Password: blockedpassword})
		assert.NoError(t, err, "authentication failed")
		_, err = handlerGrpc.token.Block(blockedUser.AccessToken.Token)
		assert.NoError(t, err, "block failed")
	})

	t.Run("FileUpload with blockedUser", func(t *testing.T) {
		// test with invalide token
		_, err = handlerGrpc.FileUpload(context.Background(),
			&grpcKeeper.UploadBinaryRequest{Name: name, Data: []byte(data),
				AccessToken: blockedUser.AccessToken})
		assert.Error(t, err, "FileUpload failed")
	})

	t.Run("FileUpload", func(t *testing.T) {
		_, err = handlerGrpc.FileUpload(context.Background(),
			&grpcKeeper.UploadBinaryRequest{Name: name, Data: []byte(data),
				AccessToken: authenticatedUser.AccessToken})
		assert.NoError(t, err, "FileUpload failed")
	})

	t.Run("FileUpload", func(t *testing.T) {
		_, err = handlerGrpc.FileUpload(context.Background(),
			&grpcKeeper.UploadBinaryRequest{Name: name, Data: []byte(data),
				AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileUpload failed")
	})

	t.Run("FileDownload with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.FileDownload(context.Background(),
			&grpcKeeper.DownloadBinaryRequest{Name: name,
				AccessToken: blockedUser.AccessToken})
		assert.Error(t, err, "FileDownload failed")
	})

	t.Run("FileDownload", func(t *testing.T) {
		_, err = handlerGrpc.FileDownload(context.Background(),
			&grpcKeeper.DownloadBinaryRequest{Name: name,
				AccessToken: authenticatedUser.AccessToken})
		assert.NoError(t, err, "FileDownload failed")
	})

	t.Run("FileRemove with blockeduser ", func(t *testing.T) {
		_, err = handlerGrpc.FileRemove(context.Background(),
			&grpcKeeper.DeleteBinaryRequest{Name: name,
				AccessToken: blockedUser.AccessToken})
		assert.Error(t, err, "FileRemove failed")
	})

	t.Run("FileRemove", func(t *testing.T) {
		_, err = handlerGrpc.FileRemove(context.Background(),
			&grpcKeeper.DeleteBinaryRequest{Name: name,
				AccessToken: authenticatedUser.AccessToken})
		assert.NoError(t, err, "FileRemove failed")
	})

	t.Run("FileDownload", func(t *testing.T) {
		_, err = handlerGrpc.FileDownload(context.Background(),
			&grpcKeeper.DownloadBinaryRequest{Name: name,
				AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileDownload failed")
	})

	t.Run("FileGetList with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.FileGetList(context.Background(),
			&grpcKeeper.GetListBinaryRequest{AccessToken: blockedUser.AccessToken})
		assert.Error(t, err, "FileGetList failed")
	})

	t.Run("FileGetList", func(t *testing.T) {
		_, err = handlerGrpc.FileGetList(context.Background(),
			&grpcKeeper.GetListBinaryRequest{AccessToken: authenticatedUser.AccessToken})
		assert.NoError(t, err, "FileGetList failed")
	})

	t.Run("create entity with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.EntityCreate(context.Background(),
			&grpcKeeper.CreateEntityRequest{
				Data: []byte(data), Metadata: string(jsonMetadata),
				AccessToken: &grpcKeeper.Token{
					Token:     blockedUser.AccessToken.Token,
					UserId:    blockedUser.AccessToken.UserId,
					CreatedAt: blockedUser.AccessToken.CreatedAt,
					EndDateAt: blockedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "create entity failed")
	})

	t.Run("create entity", func(t *testing.T) {
		jsonMetadata, _ := json.Marshal(
			model.MetadataEntity{Name: "", Description: description, Type: vars.Text.ToString()})

		_, err = handlerGrpc.EntityCreate(context.Background(),
			&grpcKeeper.CreateEntityRequest{
				Data: []byte(data), Metadata: string(jsonMetadata),
				AccessToken: &grpcKeeper.Token{
					Token:     authenticatedUser.AccessToken.Token,
					UserId:    authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt,
					EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "create entity failed")
	})

	t.Run("create entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityCreate(context.Background(),
			&grpcKeeper.CreateEntityRequest{
				Data: []byte(data), Metadata: string(jsonMetadata),
				AccessToken: &grpcKeeper.Token{
					Token:     authenticatedUser.AccessToken.Token,
					UserId:    authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt,
					EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.NoError(t, err, "create entity failed")
	})

	t.Run("create duplicate entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityCreate(context.Background(),
			&grpcKeeper.CreateEntityRequest{
				Data: []byte(data), Metadata: string(jsonMetadata),
				AccessToken: &grpcKeeper.Token{
					Token:     authenticatedUser.AccessToken.Token,
					UserId:    authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt,
					EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "create duplicate entity failed")
	})

	t.Run("update entity with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.EntityUpdate(context.Background(),
			&grpcKeeper.UpdateEntityRequest{Name: name, Data: []byte(dataUpdate), Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: blockedUser.AccessToken.Token, UserId: blockedUser.AccessToken.UserId,
					CreatedAt: blockedUser.AccessToken.CreatedAt, EndDateAt: blockedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "update entity failed")
	})

	t.Run("update entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityUpdate(context.Background(),
			&grpcKeeper.UpdateEntityRequest{Name: name, Data: []byte(dataUpdate), Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: authenticatedUser.AccessToken.Token, UserId: authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt, EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.NoError(t, err, "update entity failed")
	})

	t.Run("get list entity with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.EntityGetList(context.Background(),
			&grpcKeeper.GetListEntityRequest{Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: blockedUser.AccessToken.Token, UserId: blockedUser.AccessToken.UserId,
					CreatedAt: blockedUser.AccessToken.CreatedAt, EndDateAt: blockedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "get list failed")
	})

	t.Run("get list entity with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.EntityGetList(context.Background(),
			&grpcKeeper.GetListEntityRequest{Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: blockedUser.AccessToken.Token, UserId: blockedUser.AccessToken.UserId,
					CreatedAt: blockedUser.AccessToken.CreatedAt, EndDateAt: blockedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "get list failed")
	})

	t.Run("get list entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityGetList(context.Background(),
			&grpcKeeper.GetListEntityRequest{Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: authenticatedUser.AccessToken.Token, UserId: authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt, EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.NoError(t, err, "get list failed")
	})

	t.Run("delete entity with blockedUser", func(t *testing.T) {
		_, err = handlerGrpc.EntityDelete(context.Background(),
			&grpcKeeper.DeleteEntityRequest{Name: name, Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: blockedUser.AccessToken.Token, UserId: blockedUser.AccessToken.UserId,
					CreatedAt: blockedUser.AccessToken.CreatedAt, EndDateAt: blockedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "delete entity failed")
	})

	t.Run("delete entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityDelete(context.Background(),
			&grpcKeeper.DeleteEntityRequest{Name: name, Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: authenticatedUser.AccessToken.Token, UserId: authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt, EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.NoError(t, err, "delete entity failed")
	})

	// drop dase and test all methods again

	db.Close()

	////////////////

	t.Run("ping db", func(t *testing.T) {
		_, err = handlerGrpc.Ping(context.Background(), &grpcKeeper.PingRequest{})
		assert.Error(t, err, "failed ping db")
	})

	t.Run("registration", func(t *testing.T) {
		_, err = handlerGrpc.Registration(context.Background(), &grpcKeeper.RegistrationRequest{Username: username, Password: password})
		assert.Error(t, err, "registration failed")
	})

	t.Run("user exist", func(t *testing.T) {
		_, err = handlerGrpc.UserExist(context.Background(), &grpcKeeper.UserExistRequest{Username: username})
		assert.Error(t, err, "user exist failed")
	})

	t.Run("FileUpload", func(t *testing.T) {
		_, err = handlerGrpc.FileUpload(context.Background(),
			&grpcKeeper.UploadBinaryRequest{Name: name, Data: []byte(data),
				AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileUpload failed")
	})

	t.Run("FileDownload", func(t *testing.T) {
		_, err = handlerGrpc.FileDownload(context.Background(),
			&grpcKeeper.DownloadBinaryRequest{Name: name,
				AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileDownload failed")
	})

	t.Run("FileRemove", func(t *testing.T) {
		_, err = handlerGrpc.FileRemove(context.Background(),
			&grpcKeeper.DeleteBinaryRequest{Name: name,
				AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileRemove failed")
	})

	t.Run("FileDownload", func(t *testing.T) {
		_, err = handlerGrpc.FileDownload(context.Background(),
			&grpcKeeper.DownloadBinaryRequest{Name: name,
				AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileDownload failed")
	})

	t.Run("FileGetList", func(t *testing.T) {
		_, err = handlerGrpc.FileGetList(context.Background(),
			&grpcKeeper.GetListBinaryRequest{AccessToken: authenticatedUser.AccessToken})
		assert.Error(t, err, "FileGetList failed")
	})

	t.Run("create entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityCreate(context.Background(),
			&grpcKeeper.CreateEntityRequest{
				Data: []byte(data), Metadata: string(jsonMetadata),
				AccessToken: &grpcKeeper.Token{
					Token:     authenticatedUser.AccessToken.Token,
					UserId:    authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt,
					EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "create entity failed")
	})

	t.Run("update entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityUpdate(context.Background(),
			&grpcKeeper.UpdateEntityRequest{Name: name, Data: []byte(dataUpdate), Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: authenticatedUser.AccessToken.Token, UserId: authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt, EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "update entity failed")
	})

	t.Run("delete entity", func(t *testing.T) {
		_, err = handlerGrpc.EntityDelete(context.Background(),
			&grpcKeeper.DeleteEntityRequest{Name: name, Type: vars.Text.ToString(),
				AccessToken: &grpcKeeper.Token{Token: authenticatedUser.AccessToken.Token, UserId: authenticatedUser.AccessToken.UserId,
					CreatedAt: authenticatedUser.AccessToken.CreatedAt, EndDateAt: authenticatedUser.AccessToken.EndDateAt}})
		assert.Error(t, err, "delete entity failed")
	})

	t.Run("authentication", func(t *testing.T) {
		authenticatedUser, err = handlerGrpc.Authentication(context.Background(), &grpcKeeper.AuthenticationRequest{Username: username, Password: password})
		assert.Error(t, err, "authentication failed")
	})
}
