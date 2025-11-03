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

	// Create therapist table first (parent table for foreign key)
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS therapist (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
		)`)
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

	// Create student table (with foreign key to therapist)
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS student (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			dob DATE,
			therapist_id UUID NOT NULL REFERENCES therapist(id),
			grade INTEGER CHECK (grade >= -1 AND grade <= 12),
			iep TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		t.Fatal(err)
	}

	// Create sessions table
	_, err = pool.Exec(ctx, `
        CREATE TABLE IF NOT EXISTS sessions (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			therapist_id UUID NOT NULL,
			session_date DATE NOT NULL,
			start_time TIME,
			end_time TIME,
			notes TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
        )`)
	if err != nil {
		t.Fatal(err)
	}

	// Create theme table (parent table for foreign key)
	_, err = pool.Exec(ctx, `
		CREATE TABLE theme (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			theme_name VARCHAR(255) NOT NULL,
			month INTEGER CHECK (month >= 1 AND month <= 12),
			year INTEGER CHECK (year >= 2000 AND year <= 2500),
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
		);`)
	if err != nil {
		t.Fatal(err)
	}

	// Create SessionStudent junction table
	_, err = pool.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS session_student (
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
    );`)
	if err != nil {
		t.Fatal(err)
	}

	// Create resource table (with foreign key to theme)
	_, err = pool.Exec(ctx, `
		CREATE TABLE resource (
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
		);`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE session_resource (
    	session_id UUID,
    	resource_id UUID,
    	created_at TIMESTAMPTZ DEFAULT now(),
    	updated_at TIMESTAMPTZ DEFAULT now(),
    	PRIMARY KEY (session_id, resource_id),
    	FOREIGN KEY (session_id) REFERENCES session(id) ON DELETE CASCADE,
    	FOREIGN KEY (resource_id) REFERENCES resource(id) ON DELETE CASCADE
		);`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TYPE category AS ENUM ('visual_cue', 'verbal_cue', 'gestural_cue', 'engagement');
		CREATE TYPE response_level AS ENUM ('minimal', 'moderate', 'maximal', 'low', 'high');

		CREATE TABLE IF NOT EXISTS session_rating (
			id SERIAL PRIMARY KEY,
			session_student_id INT REFERENCES session_student(id),
			category category,
			level response_level,
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT now(),
			updated_at TIMESTAMPTZ DEFAULT now()
			);

		ALTER TABLE session_rating
		ADD CONSTRAINT unique_session_student_category 
		UNIQUE (session_student_id, category);
	`)
	if err != nil {
		t.Fatal(err)
	}

	// Create game_content Table
	_, err = pool.Exec(ctx, `
		CREATE TABLE game_content (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		category VARCHAR(255) NOT NULL CHECK (category IN ( 'sequencing', 'following_directions',
														   'wh_questions', 'true_false',
														   'concepts_sorting' )),
		level INT NOT NULL CHECK ( level >= 0 AND level <= 12 ),
		options TEXT[] NOT NULL,
		answer TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
		UNIQUE (category, level)
	);`)
	if err != nil {
		t.Fatal(err)
	}

	// Create game_result Table
	_, err = pool.Exec(ctx, `
		CREATE TABLE game_result (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		session_id UUID NOT NULL,
		student_id UUID NOT NULL,
		content_id UUID NOT NULL,
		time_taken INTEGER NOT NULL CHECK ( time_taken >= 0 ),
		completed BOOLEAN DEFAULT FALSE,
		incorrect_tries INTEGER DEFAULT 0 CHECK ( incorrect_tries >= 0 ),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now(),
	
		FOREIGN KEY (session_id, student_id) REFERENCES session_student(session_id, student_id) ON DELETE CASCADE,
		FOREIGN KEY (content_id) REFERENCES game_content(id) ON DELETE RESTRICT
	);`)
	if err != nil {
		t.Fatal(err)
	}
}
