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
		CREATE TABLE IF NOT EXISTS therapist (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
						first_name VARCHAR(100) NOT NULL,
						last_name VARCHAR(100) NOT NULL,
						email VARCHAR(255) UNIQUE NOT NULL,
						active BOOLEAN DEFAULT TRUE,
						created_at TIMESTAMPTZ DEFAULT now(),
						updated_at TIMESTAMPTZ DEFAULT now()
		)
`)
	if err != nil {
		t.Fatal(err)
	}


	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS session (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            start_datetime TIMESTAMPTZ NOT NULL,
            end_datetime TIMESTAMPTZ NOT NULL,
            therapist_id UUID NOT NULL,
            notes TEXT,
            created_at TIMESTAMPTZ DEFAULT now(),
            updated_at TIMESTAMPTZ DEFAULT now(),
            FOREIGN KEY (therapist_id) REFERENCES therapist(id) ON DELETE RESTRICT,
            CHECK (end_datetime > start_datetime)
        )
    `)
	if err != nil {
		t.Fatal(err)
	}
}
