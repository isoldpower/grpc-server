package database

import (
	"database/sql"
	"fmt"
	"golang-grpc/internal/database/types"
	"golang-grpc/internal/log"
	"strconv"

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
	config   Config
}

func NewDatabase(config Config) *Database {
	return &Database{
		instance: nil,
		config:   config,
	}
}

func (db *Database) Instantiate() types.Service {
	if db.instance != nil {
		return db.instance
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
		log.PrintError(fmt.Sprintf("Can't connect to the database at %s", connStr), err)
	}

	db.Database = databaseConnection
	db.instance = &service{
		db:       databaseConnection,
		settings: db.config,
	}

	return db.instance
}
