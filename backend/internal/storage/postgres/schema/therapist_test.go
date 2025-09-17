package schema_test

import (
	"context"
	"testing"
	"time"

	"specialstandard/internal/models"
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
	assert.Equal(t, "Matula", therapist.LastName)
	assert.Equal(t, therapistID, therapist.ID) // Optional: verify the therapist ID matches
}

func TestSessionRepository_GetTherapists(t *testing.T) {
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
	therapists, err := repo.GetTherapists(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Matula", therapists[0].LastName)
	assert.Equal(t, therapistID, therapists[0].ID) // Optional: verify the therapist ID matches
}

func TestSessionRepository_PatchTherapist(t *testing.T) {
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

	newLastName := "Matula"
	updated := &models.UpdateTherapist{
		LastName: &newLastName,
	}
	therapist, err := repo.PatchTherapist(ctx, therapistID.String(), updated)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Matula", therapist.LastName)
	assert.Equal(t, therapistID, therapist.ID) // Optional: verify the therapist ID matches
}

func TestSessionRepository_DeleteTherapist(t *testing.T) {
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

	err = repo.DeleteTherapist(ctx, therapistID.String())

	// Assert
	assert.NoError(t, err)
	//assert.Equal(t, "User " + therapistID.String() + " was deleted successfully!", mes)
}

func TestSessionRepository_CreateTherapist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewTherapistRepository(testDB.Pool)
	ctx := context.Background()

	updated := &models.CreateTherapistInput{
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:     "matulakevin91@gmai.com",
	}

	therapist, err := repo.CreateTherapist(ctx, updated)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Matula", therapist.LastName)
	assert.Equal(t, "matulakevin91@gmai.com", therapist.Email)
}
