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

func TestSessionRepository_GetTherapistByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewTherapistRepository(testDB.Pool)
	ctx := context.Background()

	// Generate a UUID for therapist_id
	therapistID := uuid.New()

	// Insert test data with UUID
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	therapist, err := repo.GetTherapistByID(ctx, therapistID.String())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Matula", therapist.Last_name)
	assert.Equal(t, therapistID, therapist.ID) // Optional: verify the therapist ID matches
}
