package events

import (
	"context"
	"crypto/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	clientConfig "github.com/xChygyNx/gophkeeper/internal/client/config"
	"github.com/xChygyNx/gophkeeper/internal/client/model"
	"github.com/xChygyNx/gophkeeper/internal/client/service/encryption"
	"github.com/xChygyNx/gophkeeper/internal/client/service/randomizer"
	grpchandler "github.com/xChygyNx/gophkeeper/internal/server/api/grpc"
	serverConfig "github.com/xChygyNx/gophkeeper/internal/server/config"
	"github.com/xChygyNx/gophkeeper/internal/server/database"
	grpcKeeper "github.com/xChygyNx/gophkeeper/internal/server/proto"
	"github.com/xChygyNx/gophkeeper/internal/server/storage"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/entity"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/file"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/token"
	"github.com/xChygyNx/gophkeeper/internal/server/storage/repositories/user"
)

func TestEvents(t *testing.T) {

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
		t.Fatalf("test containers failed: %v", err)
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
	// setting
	logger := logrus.New()
	db, err := database.New(&serverConfig.Config{DSN: databaseURI}, logger)

	if err != nil {
		t.Fatalf("db init failed: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	err = db.CreateTablesMigration("file://../../../../migrations")
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	// configs
	serverConfig := &serverConfig.Config{
		AddressGRPC:         "localhost:50051",
		DSN:                 databaseURI,
		AccessTokenLifetime: 300 * time.Second,
		DebugLevel:          logrus.DebugLevel,
		FileFolder:          "../../../../data/test_keeper",
	}

	clientConfig := &clientConfig.ConfigClient{
		GRPC:       "localhost:50051",
		DebugLevel: logrus.DebugLevel,
		FileFolder: "../../../../data/test_keeper",
	}

	// repositories
	userRepository := user.New(db)
	fileRepository := file.New(db)
	storage := storage.New("/tmp")
	entityRepository := entity.New(db)
	tokenRepository := token.New(db)

	// setup server service
	handlerGrpc := *grpchandler.NewHandler(db, serverConfig, userRepository, fileRepository, &storage,
		entityRepository, tokenRepository, logger)
	lis, err := net.Listen("tcp", serverConfig.AddressGRPC)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	t.Run("registration gophkeeper-server", func(t *testing.T) {
		grpcKeeper.RegisterGophkeeperServer(s, &handlerGrpc)
		go func() {
			if err = s.Serve(lis); err != nil {
				t.Errorf("server exited with error: %v", err)
				return
			}
		}()
	})

	connectionServer, err := grpc.Dial(serverConfig.AddressGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("connection server with error: %v", err)
	}
	grpcKeeper.NewGophkeeperClient(connectionServer)

	// setup client service
	connectionClient, err := grpc.Dial(clientConfig.GRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("connection client with error: %v", err)
	}
	grpcClient := grpcKeeper.NewGophkeeperClient(connectionClient)
	client := NewEvent(context.Background(), clientConfig, logger, grpcClient)
	client.GetConfig()

	// -- TEST DATA --
	//(name string, description string, password string, plaintext string, token model.Token)
	var accessToken model.Token = model.Token{}
	var delRow []string
	username := randomizer.RandStringRunes(10)
	password, _ := encryption.HashPassword("Password-00")
	name := randomizer.RandStringRunes(10)
	description := randomizer.RandStringRunes(10)
	plaintext := randomizer.RandStringRunes(10)
	loginUser := randomizer.RandStringRunes(10)
	passwordUser := randomizer.RandStringRunes(10)

	paymentSystem := randomizer.RandStringRunes(10)
	numberCard := "4532015112830366"
	holder := randomizer.RandStringRunes(10)
	cvc := "111"
	endDate := "01/02/2006"

	// -- TESTS --
	t.Run("ping db", func(t *testing.T) {
		_, err = client.Ping()
		assert.NoError(t, err, "failed ping db")
	})
	t.Run("registration", func(t *testing.T) {
		accessToken, err = client.Registration(username, password)
		assert.NoError(t, err, "failed registration")
	})

	t.Run("user exist", func(t *testing.T) {
		_, err = client.UserExist(username)
		assert.NoError(t, err, "failed registration")
	})

	t.Run("authentication", func(t *testing.T) {
		_, err = client.Authentication(username, password+"wrongdata")
		assert.Error(t, err, "failed authentication")
	})

	t.Run("authentication", func(t *testing.T) {
		_, err = client.Authentication(username, password)
		assert.NoError(t, err, "failed authentication")
	})

	t.Run("registration", func(t *testing.T) {
		_, err = client.Registration(username, password)
		assert.Error(t, err, "failed registration")
	})

	t.Run("FileUpload", func(t *testing.T) {
		randData := make([]byte, 255)
		_, err = rand.Read(randData)
		assert.NoError(t, err, "failed rand.Read")
		_, err = client.FileUpload(name, passwordUser, randData, accessToken)
		assert.NoError(t, err, "failed FileUpload")
	})

	t.Run("FileDownload", func(t *testing.T) {
		err = client.FileDownload(name, passwordUser, accessToken)
		assert.NoError(t, err, "failed FileDownload")
	})

	t.Run("FileRemove", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.FileRemove(delRow, accessToken)
		assert.NoError(t, err, "failed FileRemove")
	})

	t.Run("text create", func(t *testing.T) {
		err = client.TextCreate(name, description, password, plaintext, accessToken)
		assert.NoError(t, err, "failed text create")
	})

	t.Run("text update", func(t *testing.T) {
		err = client.TextUpdate(name, password, plaintext+":update", accessToken)
		assert.NoError(t, err, "failed text update")
	})

	t.Run("text delete", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.TextDelete(delRow, accessToken)
		assert.NoError(t, err, "failed text delete")
	})

	t.Run("login password create", func(t *testing.T) {
		err = client.LoginPasswordCreate(name, description, password, loginUser, passwordUser, accessToken)
		assert.NoError(t, err, "failed login password create")
	})

	t.Run("login password update", func(t *testing.T) {
		err = client.LoginPasswordUpdate(name, password, loginUser+":update", passwordUser+":update", accessToken)
		assert.NoError(t, err, "failed login password update")
	})
	t.Run("login password delete", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.LoginPasswordDelete(delRow, accessToken)
		assert.NoError(t, err, "failed login password delete")
	})
	t.Run("card create", func(t *testing.T) {
		err = client.CardCreate(name, description, password, paymentSystem, numberCard, holder, cvc, endDate, accessToken)
		assert.NoError(t, err, "failed card create")
	})
	t.Run("card update", func(t *testing.T) {
		err = client.CardUpdate(name, password, paymentSystem+":update", numberCard, holder+":update", cvc, endDate, accessToken)
		assert.NoError(t, err, "failed card update")
	})
	t.Run("card delete", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.CardDelete(delRow, accessToken)
		assert.NoError(t, err, "failed card delete")
	})
	t.Run("card delete nil", func(t *testing.T) {
		delRow = append(delRow, name)
		accessTokenTest := model.Token{}
		err = client.CardDelete(delRow, accessTokenTest)
		assert.Error(t, err, "failed card delete")
	})

	t.Run("FileUpload", func(t *testing.T) {
		randData := make([]byte, 255)
		_, err = rand.Read(randData)
		assert.NoError(t, err, "failed rand.Read")
		_, err = client.FileUpload(name, passwordUser, randData, accessToken)
		assert.NoError(t, err, "failed FileUpload")
	})

	t.Run("text create", func(t *testing.T) {
		err = client.TextCreate(name, description, password, plaintext, accessToken)
		assert.NoError(t, err, "failed text create")
	})

	t.Run("login password create", func(t *testing.T) {
		err = client.LoginPasswordCreate(name, description, password, loginUser, passwordUser, accessToken)
		assert.NoError(t, err, "failed login password create")
	})

	t.Run("card create", func(t *testing.T) {
		err = client.CardCreate(name, description, password, paymentSystem, numberCard, holder, cvc, endDate, accessToken)
		assert.NoError(t, err, "failed card create")
	})

	t.Run("card create with errors", func(t *testing.T) {
		wrongCvc := "wrongCVC"
		err = client.CardCreate(name, description, password, paymentSystem, numberCard, holder, wrongCvc, endDate, accessToken)
		assert.Error(t, err, "failed card create")
	})

	t.Run("card create with errors", func(t *testing.T) {
		wrongEndDate := "wrongEndDate"
		err = client.CardCreate(name, description, password, paymentSystem, numberCard, holder, cvc, wrongEndDate, accessToken)
		assert.Error(t, err, "failed card create")
	})

	t.Run("card update with errors", func(t *testing.T) {
		wrongCvc := "wrongCVC"
		err = client.CardUpdate(name, password, paymentSystem+":update", numberCard, holder+":update", wrongCvc, endDate, accessToken)
		assert.Error(t, err, "failed card update")
	})

	t.Run("card update with errors", func(t *testing.T) {
		wrongEndDate := "wrongEndDate"
		err = client.CardUpdate(name, password, paymentSystem+":update", numberCard, holder+":update", cvc, wrongEndDate, accessToken)
		assert.Error(t, err, "failed card update")
	})

	t.Run("Synchronization", func(t *testing.T) {
		_, _, _, _, err = client.Synchronization(password, accessToken)
		assert.NoError(t, err, "failed Synchronization")
	})

	// test with close db

	db.Close()

	t.Run("ping db", func(t *testing.T) {
		_, err = client.Ping()
		assert.Error(t, err, "failed ping db")
	})
	t.Run("registration", func(t *testing.T) {
		accessToken, err = client.Registration(username, password)
		assert.Error(t, err, "failed registration")
	})

	t.Run("FileUpload", func(t *testing.T) {
		randData := make([]byte, 255)
		_, err = rand.Read(randData)
		assert.NoError(t, err, "failed rand.Read")
		_, err = client.FileUpload(name, passwordUser, randData, accessToken)
		assert.Error(t, err, "failed FileUpload")
	})

	t.Run("FileDownload", func(t *testing.T) {
		err = client.FileDownload(name, passwordUser, accessToken)
		assert.Error(t, err, "failed FileDownload")
	})

	t.Run("FileRemove", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.FileRemove(delRow, accessToken)
		assert.Error(t, err, "failed FileRemove")
	})

	t.Run("text create", func(t *testing.T) {
		err = client.TextCreate(name, description, password, plaintext, accessToken)
		assert.Error(t, err, "failed text create")
	})

	t.Run("text update", func(t *testing.T) {
		err = client.TextUpdate(name, password, plaintext+":update", accessToken)
		assert.Error(t, err, "failed text update")
	})

	t.Run("text delete", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.TextDelete(delRow, accessToken)
		assert.Error(t, err, "failed text delete")
	})

	t.Run("login password create", func(t *testing.T) {
		err = client.LoginPasswordCreate(name, description, password, loginUser, passwordUser, accessToken)
		assert.Error(t, err, "failed login password create")
	})

	t.Run("login password update", func(t *testing.T) {
		err = client.LoginPasswordUpdate(name, password, loginUser+":update", passwordUser+":update", accessToken)
		assert.Error(t, err, "failed login password update")
	})
	t.Run("login password delete", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.LoginPasswordDelete(delRow, accessToken)
		assert.Error(t, err, "failed login password delete")
	})
	t.Run("card create", func(t *testing.T) {
		err = client.CardCreate(name, description, password, paymentSystem, numberCard, holder, cvc, endDate, accessToken)
		assert.Error(t, err, "failed card create")
	})
	t.Run("card update", func(t *testing.T) {
		err = client.CardUpdate(name, password, paymentSystem+":update", numberCard, holder+":update", cvc, endDate, accessToken)
		assert.Error(t, err, "failed card update")
	})
	t.Run("card delete", func(t *testing.T) {
		delRow = append(delRow, name)
		err = client.CardDelete(delRow, accessToken)
		assert.Error(t, err, "failed card delete")
	})

	t.Run("FileUpload", func(t *testing.T) {
		randData := make([]byte, 255)
		_, err = rand.Read(randData)
		assert.NoError(t, err, "failed rand.Read")
		_, err = client.FileUpload(name, passwordUser, randData, accessToken)
		assert.Error(t, err, "failed FileUpload")
	})

	t.Run("authentication", func(t *testing.T) {
		_, err = client.Authentication(username, password+"wrongdata")
		assert.Error(t, err, "failed authentication")
	})

	t.Run("Synchronization", func(t *testing.T) {
		_, _, _, _, err = client.Synchronization(password, accessToken)
		assert.Error(t, err, "failed Synchronization")
	})
}
