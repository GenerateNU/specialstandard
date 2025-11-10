// testutil/shared_db.go
package testutil

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	// Global shared test database - initialized once, used by all tests
	sharedTestDB *SharedTestDB
	setupOnce    sync.Once
	setupError   error
)

// SharedTestDB holds the shared test database resources
type SharedTestDB struct {
	Pool      *pgxpool.Pool
	container testcontainers.Container
	mu        sync.Mutex
}

// GetSharedTestDB returns the global shared test database
// This function is thread-safe and ensures the container is only started once
func GetSharedTestDB() (*SharedTestDB, error) {
	setupOnce.Do(func() {
		ctx := context.Background()

		log.Println("Starting shared PostgreSQL container...")

		// Start PostgreSQL container with optimizations
		pgContainer, err := postgres.Run(ctx,
			"postgres:15-alpine",
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("test"),
			postgres.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(30*time.Second),
			),
			// Add test optimizations
			testcontainers.WithEnv(map[string]string{
				"POSTGRES_INITDB_ARGS":      "-E UTF8 --auth-local=trust",
				"POSTGRES_HOST_AUTH_METHOD": "trust",
			}),
		)
		if err != nil {
			setupError = fmt.Errorf("failed to start container: %w", err)
			return
		}

		// Get connection string
		connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			setupError = fmt.Errorf("failed to get connection string: %w", err)
			return
		}

		// Configure pool for testing
		config, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			setupError = fmt.Errorf("failed to parse config: %w", err)
			return
		}

		// Optimize pool for testing
		config.MaxConns = 50 // Allow many parallel tests
		config.MinConns = 10 // Keep connections warm
		config.MaxConnLifetime = time.Hour
		config.MaxConnIdleTime = time.Minute * 30

		// Create pool
		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			setupError = fmt.Errorf("failed to create pool: %w", err)
			return
		}

		sharedTestDB = &SharedTestDB{
			Pool:      pool,
			container: pgContainer,
		}

		// Apply test optimizations to PostgreSQL
		applyTestOptimizations(pool)

		// Create all tables
		if err := createAllTables(pool); err != nil {
			setupError = fmt.Errorf("failed to create tables: %w", err)
			return
		}

		log.Println("Shared test database ready!")
	})

	return sharedTestDB, setupError
}

// applyTestOptimizations applies PostgreSQL settings for faster tests
func applyTestOptimizations(pool *pgxpool.Pool) {
	ctx := context.Background()

	// ONLY for tests - makes writes much faster but less safe
	optimizations := []string{
		"ALTER SYSTEM SET fsync = off",
		"ALTER SYSTEM SET synchronous_commit = off",
		"ALTER SYSTEM SET full_page_writes = off",
		"ALTER SYSTEM SET checkpoint_segments = 100",
		"ALTER SYSTEM SET checkpoint_completion_target = 0.9",
		"ALTER SYSTEM SET wal_buffers = '64MB'",
		"ALTER SYSTEM SET shared_buffers = '256MB'",
		"SELECT pg_reload_conf()",
	}

	for _, sql := range optimizations {
		if _, err := pool.Exec(ctx, sql); err != nil {
			log.Printf("Warning: Failed to apply optimization %q: %v", sql, err)
		}
	}
}

// createAllTables creates all database tables
func createAllTables(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Use a transaction for atomic table creation
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Enable extensions
	if _, err := tx.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	`); err != nil {
		return fmt.Errorf("failed to create extensions: %w", err)
	}

	// Create tables in dependency order
	tableDefinitions := []string{
		`CREATE TABLE IF NOT EXISTS district (
  			id SERIAL PRIMARY KEY,
 	 		name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS school (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			district_id INTEGER NOT NULL REFERENCES district(id),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS therapist (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now(),
			schools INTEGER[],
			district_id INTEGER NOT NULL REFERENCES district(id)
		)`,

		`CREATE TABLE IF NOT EXISTS theme (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			theme_name VARCHAR(255) NOT NULL,
			month INTEGER CHECK (month >= 1 AND month <= 12),
			year INTEGER CHECK (year >= 2000 AND year <= 2500),
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
		)`,

		`CREATE TABLE IF NOT EXISTS student (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			dob DATE,
			therapist_id UUID NOT NULL REFERENCES therapist(id),
			school_id INTEGER NOT NULL REFERENCES school(id),
			grade INTEGER CHECK (grade >= -1 AND grade <= 12),
			iep TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS session (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			start_datetime TIMESTAMPTZ NOT NULL,
			end_datetime TIMESTAMPTZ NOT NULL,
			therapist_id UUID NOT NULL,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now(),
			FOREIGN KEY (therapist_id) REFERENCES therapist(id) ON DELETE RESTRICT,
			CHECK (end_datetime > start_datetime)
		)`,

		`CREATE TABLE IF NOT EXISTS session_student (
			id SERIAL PRIMARY KEY,
			session_id UUID,
			student_id UUID,
			present BOOLEAN DEFAULT TRUE,
			notes TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now(),
			FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
			FOREIGN KEY (student_id) REFERENCES student(id) ON DELETE CASCADE,
			UNIQUE (session_id, student_id)
		)`,

		`CREATE TABLE IF NOT EXISTS resource (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			theme_id UUID NOT NULL,
			grade_level INTEGER CHECK (grade_level >= 0 AND grade_level <= 12),
			date DATE,
			type VARCHAR(50),
			title VARCHAR(100),
			category VARCHAR(100),
			content TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now(),
			FOREIGN KEY (theme_id) REFERENCES theme(id) ON DELETE RESTRICT
		)`,

		`CREATE TABLE IF NOT EXISTS session_resource (
			session_id UUID,
			resource_id UUID,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now(),
			PRIMARY KEY (session_id, resource_id),
			FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
			FOREIGN KEY (resource_id) REFERENCES resource(id) ON DELETE CASCADE
		)`,
	}

	// Execute non-enum table creations
	for _, tableDef := range tableDefinitions {
		if _, err := tx.Exec(ctx, tableDef); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Create enums and dependent tables separately to handle "type already exists" gracefully
	if _, err := tx.Exec(ctx, `
		DO $$ 
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'category') THEN
				CREATE TYPE category AS ENUM ('visual_cue', 'verbal_cue', 'gestural_cue', 'engagement');
			END IF;
			
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'response_level') THEN
				CREATE TYPE response_level AS ENUM ('minimal', 'moderate', 'maximal', 'low', 'high');
			END IF;
		END$$;

		CREATE TABLE IF NOT EXISTS session_rating (
			id SERIAL PRIMARY KEY,
			session_student_id INT REFERENCES session_student(id),
			category category,
			level response_level,
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
		);

		-- Add constraint if it doesn't exist
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'unique_session_student_category'
			) THEN
				ALTER TABLE session_rating
				ADD CONSTRAINT unique_session_student_category 
				UNIQUE (session_student_id, category);
			END IF;
		END$$;
	`); err != nil {
		return fmt.Errorf("failed to create enums and rating table: %w", err)
	}

	return tx.Commit(ctx)
}

// CleanupTestData truncates all tables efficiently
// This is much faster than DELETE and resets auto-increment counters
func (db *SharedTestDB) CleanupTestData(t testing.TB) {
	t.Helper()

	// Ensure cleanup is thread-safe
	db.mu.Lock()
	defer db.mu.Unlock()

	ctx := context.Background()

	// Single query to truncate all tables
	// RESTART IDENTITY resets sequences
	// CASCADE handles foreign key dependencies
	_, err := db.Pool.Exec(ctx, `
		TRUNCATE TABLE 
			session_rating,
			session_resource,
			session_student,
			resource,
			student,
			session,
			theme,
			therapist,
			school,
			district
		RESTART IDENTITY CASCADE
	`)

	if err != nil {
		t.Fatalf("Failed to cleanup test data: %v", err)
	}
}

// Shutdown closes the shared test database
// Call this in TestMain after all tests complete
func Shutdown() {
	if sharedTestDB != nil {
		sharedTestDB.Pool.Close() // Fixed: was sharedtestDB.Close()

		if sharedTestDB.container != nil {
			ctx := context.Background()
			if err := sharedTestDB.container.Terminate(ctx); err != nil {
				log.Printf("Warning: Failed to terminate container: %v", err)
			}
		}
	}
}

// SetupTestWithCleanup is a helper that gets the shared DB and ensures cleanup
func SetupTestWithCleanup(t testing.TB) *pgxpool.Pool {
	t.Helper()

	db, err := GetSharedTestDB()
	if err != nil {
		t.Fatalf("Failed to get shared test database: %v", err)
	}

	// Clean before test (not after, so failed tests leave data for debugging)
	db.CleanupTestData(t)

	return db.Pool
}
