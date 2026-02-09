package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sindrema/teslamap/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser, cfg.DatabasePass, cfg.DatabaseHost, cfg.DatabasePort, cfg.DatabaseName,
	)
	pool, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	pool.SetMaxOpenConns(5)
	if err := pool.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	return pool
}
