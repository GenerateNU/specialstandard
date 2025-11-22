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
	// First, ensure district exists
	districtID := 1
	_, err := db.Exec(ctx, `
		INSERT INTO district (id, name, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, districtID, "Test District")
	assert.NoError(t, err)

	// Then, ensure school exists
	schoolID := 1
	_, err = db.Exec(ctx, `
		INSERT INTO school (id, name, district_id, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, schoolID, "Test School", districtID)
	assert.NoError(t, err)

	// Now create the therapist with required fields
	therapistID := uuid.New()
	email := fmt.Sprintf("therapist_%s@example.com", therapistID.String()[:8])
	_, err = db.Exec(ctx, `
		INSERT INTO therapist (id, first_name, last_name, email, schools, district_id, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, therapistID, "Test", "Therapist", email, []int{schoolID}, districtID, true, time.Now(), time.Now())
	assert.NoError(t, err)

	return therapistID
}

func TestSessionResourceRepository_PostSessionResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewSessionResourceRepository(testDB)
	ctx := context.Background()

	themeID := CreateTestTheme(t, testDB, ctx)
	therapistID := CreateTestTherapist(t, testDB, ctx)

	// Create test resource
	resourceID := uuid.New()
	_, err := testDB.Exec(ctx, `
		INSERT INTO resource (id, theme_id, grade_level, week, type, title, category, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, resourceID, themeID, 5, 1, "worksheet", "Animal Worksheet", "speech", "Animal recognition", time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	sessionParentID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)

	startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
	_, err = testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
		INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, sessionID, "Test Session", startTime, endTime, "Test Session", time.Now(), time.Now(), sessionParentID)
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
		sessionParentID := uuid.New()
		startTime := time.Now()
		endTime := startTime.Add(time.Hour)

		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
		_, err = testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
			INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, newSessionID, "Test Session", startTime, endTime, "Test Session 2", time.Now(), time.Now(), sessionParentID)
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
	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewSessionResourceRepository(testDB)
	ctx := context.Background()

	themeID := CreateTestTheme(t, testDB, ctx)
	therapistID := CreateTestTherapist(t, testDB, ctx)

	// Create test resource
	resourceID := uuid.New()
	_, err := testDB.Exec(ctx, `
		INSERT INTO resource (id, theme_id, grade_level, week, type, title, category, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, resourceID, themeID, 5, 1, "worksheet", "Animal Worksheet", "speech", "Animal recognition", time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	sessionParentID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)

	startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
	_, err = testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
	assert.NoError(t, err)
	_, err = testDB.Exec(ctx, `
		INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, sessionID, "Test Session", startTime, endTime, "Test Session", time.Now(), time.Now(), sessionParentID)
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
		err = testDB.QueryRow(ctx, `SELECT COUNT(*) FROM session_resource WHERE session_id = $1 AND resource_id = $2`, sessionID, resourceID).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count, "Session resource should have been deleted")
	})
}

func TestSessionResourceRepository_GetResourcesBySessionID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	var err error

	repo := schema.NewSessionResourceRepository(testDB)
	ctx := context.Background()

	themeID := CreateTestTheme(t, testDB, ctx)
	therapistID := CreateTestTherapist(t, testDB, ctx)

	t.Run("session with one resource", func(t *testing.T) {
		// Create test resource
		resourceID := uuid.New()
		_, err = testDB.Exec(ctx, `
			INSERT INTO resource (id, theme_id, grade_level, week, type, title, category, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, resourceID, themeID, 5, 1, "worksheet", "Animal Worksheet", "speech", "Animal recognition", time.Now(), time.Now())
		assert.NoError(t, err)

		sessionID := uuid.New()
		sessionParentID := uuid.New()
		startTime := time.Now()
		endTime := startTime.Add(time.Hour)

		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
		_, err = testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
			INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, sessionID, "Test Session", startTime, endTime, "Test Session", time.Now(), time.Now(), sessionParentID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
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
		startTime := time.Now()
		endTime := startTime.Add(time.Hour)

		sessionParentID := uuid.New()
		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
		_, err := testDB.Exec(ctx, `
				INSERT INTO session_parent (id, start_date, end_date, therapist_id)
				VALUES ($1, $2, $3, $4)
		`, sessionParentID, startDate, endDate, therapistID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
			INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, sessionID, "Test Session", startTime, endTime, "Multi Resource Session", time.Now(), time.Now(), sessionParentID)
		assert.NoError(t, err)

		// Create multiple resources
		resourceID1 := uuid.New()
		_, err = testDB.Exec(ctx, `
			INSERT INTO resource (id, theme_id, grade_level, week, type, title, category, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, resourceID1, themeID, 5, 1, "worksheet", "Math Worksheet", "math", "Basic arithmetic", time.Now(), time.Now())
		assert.NoError(t, err)

		resourceID2 := uuid.New()
		_, err = testDB.Exec(ctx, `
			INSERT INTO resource (id, theme_id, grade_level, week, type, title, category, content, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, resourceID2, themeID, 5, 1, "activity", "Reading Activity", "language", "Comprehension exercise", time.Now(), time.Now())
		assert.NoError(t, err)

		// Link resources to session
		_, err = testDB.Exec(ctx, `
			INSERT INTO session_resource (session_id, resource_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4)
		`, sessionID, resourceID1, time.Now(), time.Now())
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
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
		sessionParentID := uuid.New()
		startTime := time.Now()
		endTime := startTime.Add(time.Hour)

		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
		_, err = testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
			INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, sessionID, "Test Session", startTime, endTime, "Empty Session", time.Now(), time.Now(), sessionParentID)
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

		date := time.Now()

		startTime := date.Add(-1 * time.Hour)
		endTime := date

		sessionParentID := uuid.New()
		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
		_, err := testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
			INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, sessionID, "Test Session", startTime, endTime, "Multi Resource Session", time.Now(), time.Now(), sessionParentID)
		assert.NoError(t, err)

		// Create multiple resources
		// for i := 1; i <= 12; i++ {
		// 	resourceID := uuid.New()
		// 	_, err := testDB.Exec(ctx, `
		//             INSERT INTO resource (id, theme_id, grade_level, week, type, title, category, content, created_at, updated_at)
		//             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		// 		resourceID, themeID, i, 1, "activity", "Reading Activity", "language", "Comprehension", time.Now(), time.Now())
		// 	assert.NoError(t, err)

		// 	_, err = testDB.Exec(ctx, `
		//             INSERT INTO session_resource (session_id, resource_id, created_at, updated_at)
		//             VALUES ($1, $2, $3, $4)
		//             `, sessionID, resourceID, time.Now(), time.Now())
		// 	assert.NoError(t, err)
		// }

		//resources, err := repo.GetResourcesBySessionID(ctx, sessionID, utils.NewPagination())
		// assert.NoError(t, err)
		// //		assert.Len(t, resources, 12, "Expected 12 as per default pagination")

		// resources, err = repo.GetResourcesBySessionID(ctx, sessionID, utils.Pagination{
		// 	Page:  2,
		// 	Limit: 13,
		// })
		// assert.NoError(t, err)
		// assert.Len(t, resources, 0, "Expected Length 0 on last page")
	})
}
