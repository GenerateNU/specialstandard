package testutil

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	Pool      *pgxpool.Pool
	container testcontainers.Container
}

func SetupTestDB(t testing.TB) *TestDB {
	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	// Connect to database
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}

	// Create tables
	createTables(t, pool)

	return &TestDB{
		Pool:      pool,
		container: pgContainer,
	}
}

func (db *TestDB) Cleanup() {
	db.Pool.Close()
	if err := db.container.Terminate(context.Background()); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
}

func createTables(t testing.TB, pool *pgxpool.Pool) {
	ctx := context.Background()

	// Enable UUID extension first
	_, err := pool.Exec(ctx, `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		t.Fatal(err)
	}

	// Create tables here
	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS sessions (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
						therapist_id UUID NOT NULL ,
						session_date DATE NOT NULL,
						start_time TIME,
						end_time TIME,
						notes TEXT,
						created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
						updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )
    `)
	if err != nil {
		t.Fatal(err)
	}
}
