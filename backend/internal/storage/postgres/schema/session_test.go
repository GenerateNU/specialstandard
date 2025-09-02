package schema_test

import (
	"context"
	"testing"
	"time"

	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

	// Generate a UUID for therapist_id
	therapistID := uuid.New()

	// Insert test data with UUID
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO sessions (therapist_id, session_date, start_time, end_time, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, therapistID, time.Now(), "10:00", "11:00", "Test session")
	assert.NoError(t, err)

	// Test
	sessions, err := repo.GetSessions(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, "Test session", *sessions[0].Notes)
	assert.Equal(t, therapistID, sessions[0].TherapistID) // Optional: verify the therapist ID matches
}
