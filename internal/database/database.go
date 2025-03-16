package database

import (
	"context"
	"database/sql"
	"fmt"
	"golang-grpc/internal/database/types"
	"golang-grpc/internal/log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Schema   string
}

type Database struct {
	Database *sql.DB
	instance *service
	config   *Config
}

func NewDatabase(config *Config) *Database {
	return &Database{
		instance: nil,
		config:   config,
	}
}

func (db *Database) pingConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Database.PingContext(ctx); err != nil {
		db.Database.Close()
		return err
	}

	return nil
}

func (db *Database) Instantiate() (types.Service, error) {
	if db.instance != nil {
		return db.instance, nil
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		db.config.Username,
		db.config.Password,
		db.config.Host,
		strconv.Itoa(db.config.Port),
		db.config.Database,
		db.config.Schema,
	)

	databaseConnection, err := sql.Open("postgres", connStr)
	if err != nil {
		log.PrintError(fmt.Sprintf("Failed to instantiate connection to database at %s", connStr), err)
	}
	db.Database = databaseConnection
	db.instance = &service{
		db:       databaseConnection,
		settings: db.config,
	}

	if connectionErr := db.pingConnection(); connectionErr != nil {
		log.PrintError(fmt.Sprintf("Can't connect to the database at %s", connStr), connectionErr)
		return nil, connectionErr
	}

	return db.instance, nil
}
