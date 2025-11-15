package postgres

import (
	"specialstandard/internal/config"
	"specialstandard/internal/storage"

	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Establishes a sustained connection to the PostgreSQL database using pooling.
func ConnectDatabase(ctx context.Context, config config.DB) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(config.Connection())
	//dbConfig, err := pgxpool.ParseConfig("postgresql://postgres:postgres@host.docker.internal:54322/postgres")
	// dbConfig, err := pgxpool.ParseConfig("postgresql://postgres:postgres@127.0.0.1:54322/postgres")

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}
	// Disable prepared statements to avoid conflicts during hot reload in development
	// This prevents "prepared statement already exists" errors when connections are reused
	dbConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	log.Print("Connected to database!")
	return conn, nil
}

func NewRepository(ctx context.Context, config config.DB) *storage.Repository {
	db, err := ConnectDatabase(ctx, config)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	return storage.NewRepository(db)
}
