package schema_test

import (
	"context"
	"testing"
	"time"

	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionRepository_GetSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionRepository(testDB.Pool)
	ctx := context.Background()

	// Create a test therapist first (required for foreign key)
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email)
        VALUES ($1, $2, $3, $4)
    `, therapistID, "John", "Doe", "john.doe@example.com")
	assert.NoError(t, err)

	// Insert test session data using new schema
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session (therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4)
    `, therapistID, startTime, endTime, "Test session")
	assert.NoError(t, err)

	// Test
	sessions, err := repo.GetSessions(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, "Test session", *sessions[0].Notes)
	assert.Equal(t, therapistID, sessions[0].TherapistID)
	assert.True(t, sessions[0].EndDatetime.After(sessions[0].StartDatetime))
}

func TestSessionRepository_GetSessionByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionRepository(testDB.Pool)
	ctx := context.Background()

	// Create a test therapist first
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email)
        VALUES ($1, $2, $3, $4)
    `, therapistID, "Jane", "Smith", "jane.smith@example.com")
	assert.NoError(t, err)

	// Insert test session and capture the generated ID
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Get by ID test session")
	assert.NoError(t, err)

	// Test
	session, err := repo.GetSessionByID(ctx, sessionID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, sessionID, session.ID)
	assert.Equal(t, therapistID, session.TherapistID)
	assert.Equal(t, "Get by ID test session", *session.Notes)

	// Test not found
	nonExistentID := uuid.New()
	session, err = repo.GetSessionByID(ctx, nonExistentID.String())
	assert.Error(t, err)
	assert.Nil(t, session)
}

func TestSessionRepository_DeleteSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB tests in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionRepository(testDB.Pool)
	ctx := context.Background()

	// 404 NOT FOUND TEST
	badID := uuid.New()
	message, err := repo.DeleteSessions(ctx, badID)
	assert.Equal(t, err, nil)
	assert.Equal(t, message, "")

	// SUCCESS TEST - Creation of valid Therapist first.
	therapistID := uuid.New()
	_, err = testDB.Pool.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email)
             VALUES ($1, $2, $3, $4)`,
		therapistID, "Doctor", "Suess", "dr.guesswho.suess@drdr.com")
	assert.NoError(t, err)

	// Inserting test session into the DB.
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx,
		`INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
             VALUES ($1, $2, $3, $4, $5)`,
		sessionID, therapistID, startTime, endTime, "Inserting into session for test")
	assert.NoError(t, err)

	msg, err := repo.DeleteSessions(ctx, sessionID)
	require.NoError(t, err)
	assert.Equal(t, "Deleted the Session Successfully!", msg)
}
