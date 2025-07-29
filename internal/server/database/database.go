package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/xChygyNx/gophkeeper/internal/server/config"
)

type DB struct {
	Pool *sql.DB
}

// New - open connection with database
func New(config *config.Config, log *logrus.Logger) (*DB, error) {
	pool, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, err
	}

	log.Info("Connect to DB")

	ctx, cnl := context.WithTimeout(context.Background(), 10*time.Second)
	defer cnl()

	if err := pool.PingContext(ctx); err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

// Close - closes the database connection
func (db DB) Close() error {
	return db.Pool.Close()
}

// Ping - checks the database connection
func (db DB) Ping() error {
	if err := db.Pool.Ping(); err != nil {
		//db.log.Error(err)
		return err
	}
	return nil
}

// CreateTablesMigration - creates database tables using migrations
func (db DB) CreateTablesMigration(migrationSource string) error {

	driver, err := postgres.WithInstance(db.Pool, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		"postgres", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		fmt.Println("Nothing to migrate")
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
