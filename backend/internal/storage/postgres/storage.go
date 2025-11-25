package postgres

import (
	"context"
	"log"
	"time"

	"specialstandard/internal/config"
	"specialstandard/internal/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Establishes a sustained connection to the PostgreSQL database using pooling.
func ConnectDatabase(ctx context.Context, dbConfig config.DB) (*pgxpool.Pool, error) {
	log.Printf("Pool config - MaxConns: %d, MinConns: %d, MaxLifetime: %d",
		dbConfig.MaxOpenConns, dbConfig.MaxIdleConns, dbConfig.ConnMaxLifetime)

	poolConfig, err := pgxpool.ParseConfig(dbConfig.Connection())
	if err != nil {
		log.Fatalf("Failed to parse database config: %v", err)
		return nil, err
	}

	// Configure connection pool
	maxConns := dbConfig.MaxOpenConns
	if maxConns == 0 {
		maxConns = 10
	}
	poolConfig.MaxConns = int32(maxConns)

	// Set min connections to keep warm
	minConns := dbConfig.MaxIdleConns
	if minConns == 0 {
		minConns = 3
	}
	poolConfig.MinConns = int32(minConns)

	maxConnIdleTime := dbConfig.MaxConnIdleTime // 2 minutes
	if maxConnIdleTime == 0 {
		maxConnIdleTime = 10 * 60 // 5 minutes
	}
	poolConfig.MaxConnIdleTime = time.Duration(maxConnIdleTime)

	// Set max connection lifetime (prevents stale connections)
	maxLifetime := dbConfig.ConnMaxLifetime
	if maxLifetime == 0 {
		maxLifetime = 5 * 60 // 5 minutes
	}
	poolConfig.MaxConnLifetime = time.Duration(maxLifetime) * time.Second

	// Set idle timeout (closes idle connections after this duration)
	poolConfig.MaxConnIdleTime = 2 * time.Minute

	// Disable prepared statements to avoid conflicts during hot reload in development
	// This prevents "prepared statement already exists" errors when connections are reused
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Printf("New connection established")
		return nil
	}

	// Create connection pool
	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Printf("Failed to create connection pool: %v", err)
		return nil, err
	}

	// Test the connection with timeout
	pingCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	err = conn.Ping(pingCtx)
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		conn.Close()
		return nil, err
	}

	log.Print("Connected to database!")
	return conn, nil
}

func NewRepository(ctx context.Context, dbConfig config.DB) *storage.Repository {
	db, err := ConnectDatabase(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	return storage.NewRepository(db)
}
