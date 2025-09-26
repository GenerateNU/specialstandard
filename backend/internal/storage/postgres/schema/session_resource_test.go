package schema_test

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"
	"specialstandard/internal/utils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func CreateTestTheme(t *testing.T, db *pgxpool.Pool, ctx context.Context) uuid.UUID {
	themeID := uuid.New()
	_, err := db.Exec(ctx, `
		INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, themeID, "Test Theme", 1, 2024, time.Now(), time.Now())
	assert.NoError(t, err)
	return themeID
}

func CreateTestTherapist(t *testing.T, db *pgxpool.Pool, ctx context.Context) uuid.UUID {
	therapistID := uuid.New()
	email := fmt.Sprintf("therapist_%s@example.com", therapistID.String()[:8])
	_, err := db.Exec(ctx, `
		INSERT INTO therapist (id, first_name, last_name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, therapistID, "Test", "Therapist", email, time.Now(), time.Now())
	assert.NoError(t, err)
	return therapistID
}

func TestSessionResourceRepository_PostSessionResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionResourceRepository(testDB.Pool)
	ctx := context.Background()

	themeID := CreateTestTheme(t, testDB.Pool, ctx)
	therapistID := CreateTestTherapist(t, testDB.Pool, ctx)

	// Create test resource
	resourceID := uuid.New()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, resourceID, themeID, "5th Grade", testDate, "worksheet", "Animal Worksheet", "speech", "Animal recognition", time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, sessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Test Session", time.Now(), time.Now())
	assert.NoError(t, err)

	t.Run("successful creation", func(t *testing.T) {
		newSessionResource, err := repo.PostSessionResource(ctx, models.CreateSessionResource{
			SessionID:  sessionID,
			ResourceID: resourceID,
		})
		if assert.NoError(t, err) {
			assert.Equal(t, sessionID, newSessionResource.SessionID)
			assert.Equal(t, resourceID, newSessionResource.ResourceID)
			assert.NotNil(t, newSessionResource)
		}
	})

	t.Run("duplicate session resource - should fail", func(t *testing.T) {
		// Try to create the same relationship again
		_, err := repo.PostSessionResource(ctx, models.CreateSessionResource{
			SessionID:  sessionID,
			ResourceID: resourceID,
		})
		assert.Error(t, err, "Should fail on duplicate session-resource relationship")
	})

	t.Run("non-existent session - should fail", func(t *testing.T) {
		nonExistentSessionID := uuid.New()
		_, err := repo.PostSessionResource(ctx, models.CreateSessionResource{
			SessionID:  nonExistentSessionID,
			ResourceID: resourceID,
		})
		assert.Error(t, err, "Should fail with non-existent session ID")
	})

	t.Run("non-existent resource - should fail", func(t *testing.T) {
		// Create new session for this test
		newSessionID := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, newSessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Test Session 2", time.Now(), time.Now())
		assert.NoError(t, err)

		nonExistentResourceID := uuid.New()
		_, err := repo.PostSessionResource(ctx, models.CreateSessionResource{
			SessionID:  newSessionID,
			ResourceID: nonExistentResourceID,
		})
		assert.Error(t, err, "Should fail with non-existent resource ID")
	})
}

func TestSessionResourceRepository_DeleteSessionResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionResourceRepository(testDB.Pool)
	ctx := context.Background()

	themeID := CreateTestTheme(t, testDB.Pool, ctx)
	therapistID := CreateTestTherapist(t, testDB.Pool, ctx)

	// Create test resource
	resourceID := uuid.New()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, resourceID, themeID, "5th Grade", testDate, "worksheet", "Animal Worksheet", "speech", "Animal recognition", time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, sessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Test Session", time.Now(), time.Now())
	assert.NoError(t, err)

	t.Run("successful deletion", func(t *testing.T) {
		// Create the session resource to be deleted
		_, err = repo.PostSessionResource(ctx, models.CreateSessionResource{
			SessionID:  sessionID,
			ResourceID: resourceID,
		})
		assert.NoError(t, err)

		// Delete the session resource
		err = repo.DeleteSessionResource(ctx, models.DeleteSessionResource{
			SessionID:  sessionID,
			ResourceID: resourceID,
		})
		assert.NoError(t, err)

		var count int
		err = testDB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM session_resource WHERE session_id = $1 AND resource_id = $2`, sessionID, resourceID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Session resource should have been deleted")
	})
}

func TestSessionResourceRepository_GetResourcesBySessionID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()
	var err error

	repo := schema.NewSessionResourceRepository(testDB.Pool)
	ctx := context.Background()

	themeID := CreateTestTheme(t, testDB.Pool, ctx)
	therapistID := CreateTestTherapist(t, testDB.Pool, ctx)

	t.Run("session with one resource", func(t *testing.T) {
		// Create test resource
		resourceID := uuid.New()
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, resourceID, themeID, "5th Grade", testDate, "worksheet", "Animal Worksheet", "speech", "Animal recognition", time.Now(), time.Now())
		assert.NoError(t, err)

		sessionID := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, sessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Test Session", time.Now(), time.Now())
		assert.NoError(t, err)

		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session_resource (session_id, resource_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
		`, sessionID, resourceID, time.Now(), time.Now())
		assert.NoError(t, err)

		resources, err := repo.GetResourcesBySessionID(ctx, sessionID, utils.NewPagination())

		// Assert
		assert.NoError(t, err)
		if assert.NotEmpty(t, resources, "Expected exactly 1 resource") {
			assert.Len(t, resources, 1)
			assert.Equal(t, "Animal Worksheet", *resources[0].Title)
			assert.Equal(t, resourceID, resources[0].ID)
		}
	})

	t.Run("session with multiple resources", func(t *testing.T) {
		// Create session with new ID
		sessionID := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, sessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Multi Resource Session", time.Now(), time.Now())
		assert.NoError(t, err)

		// Create multiple resources
		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		resourceID1 := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, resourceID1, themeID, "5th Grade", testDate, "worksheet", "Math Worksheet", "math", "Basic arithmetic", time.Now(), time.Now())
		assert.NoError(t, err)

		resourceID2 := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, resourceID2, themeID, "5th Grade", testDate, "activity", "Reading Activity", "language", "Comprehension exercise", time.Now(), time.Now())
		assert.NoError(t, err)

		// Link resources to session
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session_resource (session_id, resource_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
		`, sessionID, resourceID1, time.Now(), time.Now())
		assert.NoError(t, err)

		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session_resource (session_id, resource_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
		`, sessionID, resourceID2, time.Now(), time.Now())
		assert.NoError(t, err)

		resources, err := repo.GetResourcesBySessionID(ctx, sessionID, utils.NewPagination())
		assert.NoError(t, err)
		assert.Len(t, resources, 2, "Expected exactly 2 resources")
	})

	t.Run("session with no resources - returns empty array", func(t *testing.T) {
		sessionID := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, sessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Empty Session", time.Now(), time.Now())
		assert.NoError(t, err)

		resources, err := repo.GetResourcesBySessionID(ctx, sessionID, utils.NewPagination())
		assert.NoError(t, err)
		assert.NotNil(t, resources, "Should return non-nil slice")
		assert.Empty(t, resources, "Should return empty array for session with no resources")
		assert.Len(t, resources, 0)
	})

	t.Run("non-existent session - returns empty array", func(t *testing.T) {
		nonExistentSessionID := uuid.New()
		resources, err := repo.GetResourcesBySessionID(ctx, nonExistentSessionID, utils.NewPagination())
		assert.NoError(t, err)
		assert.NotNil(t, resources, "Should return non-nil slice")
		assert.Empty(t, resources, "Should return empty array for non-existent session")
		assert.Len(t, resources, 0)
	})

	t.Run("More Test Cases for Pagination", func(t *testing.T) {
		// Create session with new ID
		sessionID := uuid.New()
		_, err = testDB.Pool.Exec(ctx, `
			INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, sessionID, therapistID, time.Now(), time.Now().Add(time.Hour), "Multi Resource Session", time.Now(), time.Now())
		assert.NoError(t, err)

		testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		// Create multiple resources
		for i := 1; i <= 15; i++ {
			resourceID := uuid.New()
			_, err := testDB.Pool.Exec(ctx, `
                INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
				resourceID, themeID, fmt.Sprintf("Grade = %d", i), testDate, "activity", "Reading Activity", "language", "Comprehension", time.Now(), time.Now())
			assert.NoError(t, err)

			_, err = testDB.Pool.Exec(ctx, `
                INSERT INTO session_resource (session_id, resource_id, created_at, updated_at)
                VALUES ($1, $2, $3, $4)
                `, sessionID, resourceID, time.Now(), time.Now())
			assert.NoError(t, err)
		}

		resources, err := repo.GetResourcesBySessionID(ctx, sessionID, utils.NewPagination())
		assert.NoError(t, err)
		assert.Len(t, resources, 10, "Expected 10 as per default pagination")

		resources, err = repo.GetResourcesBySessionID(ctx, sessionID, utils.Pagination{
			Page:  2,
			Limit: 13,
		})
		assert.NoError(t, err)
		assert.Len(t, resources, 2, "Expected Length 2 on last page")
	})
}
