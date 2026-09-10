package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func OpenFromEnvironment() (*sql.DB, error) {
	cfg := DBConfig{
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	}
	if cfg.Port == "" {
		cfg.Port = "3306"
	}
	if cfg.Username == "" || cfg.Host == "" || cfg.Database == "" {
		return nil, fmt.Errorf("DB_USER, DB_HOST, and DB_NAME are required")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open MySQL connection: %w", err)
	}
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("connect to MySQL: %w", err)
	}
	return database, nil
}
