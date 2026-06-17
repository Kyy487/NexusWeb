package database

import (
	"context"
	"fmt"
	"log"

	"nexusweb-market/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping database: ", err)
	}

	log.Println("database connected successfully")
	return db
}