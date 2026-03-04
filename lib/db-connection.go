package lib

import (
	"database/sql"
	"fmt"
	"log"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func DbConnection(cfg DBConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao configurar o banco: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao conectar no MariaDB: %v", err)
	}

	log.Println("✅ Conectado ao MariaDB com sucesso!")

	return db
}