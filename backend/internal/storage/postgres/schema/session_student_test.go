package schema_test

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions
func ptrBool(b bool) *bool {
	return &b
}

// CreateTestDistrict creates a test district and returns its ID
func CreateTestDistrict(t *testing.T, db *pgxpool.Pool, ctx context.Context) int {
	districtID := 1
	_, err := db.Exec(ctx, `
		INSERT INTO district (id, name, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, districtID, "Test District")
	assert.NoError(t, err)
	return districtID
}

// CreateTestSchool creates a test school and returns its ID
func CreateTestSchool(t *testing.T, db *pgxpool.Pool, ctx context.Context, districtID int) int {
	schoolID := 1
	_, err := db.Exec(ctx, `
		INSERT INTO school (id, name, district_id, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, schoolID, "Test School", districtID)
	assert.NoError(t, err)
	return schoolID
}

// CreateTestStudent creates a test student with required school
func CreateTestStudent(t *testing.T, db *pgxpool.Pool, ctx context.Context, therapistID uuid.UUID, name string) uuid.UUID {
	// Ensure school exists for student
	districtID := CreateTestDistrict(t, db, ctx)
	schoolID := CreateTestSchool(t, db, ctx, districtID)

	studentID := uuid.New()
	_, err := db.Exec(ctx, `
		INSERT INTO student (id, first_name, last_name, therapist_id, school_id, grade, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, studentID, name, "Student", therapistID, schoolID, 5)
	assert.NoError(t, err)

	return studentID
}

// CreateTestSession creates a test session
func CreateTestSession(t *testing.T, db *pgxpool.Pool, ctx context.Context, therapistID uuid.UUID, notes string) uuid.UUID {
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err := db.Exec(ctx, `
		INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`, sessionID, therapistID, startTime, endTime, notes)
	assert.NoError(t, err)

	return sessionID
}

func TestSessionStudentRepository_CreateSessionStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionStudentRepository(testDB)
	ctx := context.Background()

	// Create test data using helper functions
	therapistID := CreateTestTherapist(t, testDB, ctx)
	sessionID := CreateTestSession(t, testDB, ctx, therapistID, "Test session for session-student")
	studentID := CreateTestStudent(t, testDB, ctx, therapistID, "Test")

	// Test successful creation
	input := &models.CreateSessionStudentInput{
		SessionIDs: []uuid.UUID{sessionID},
		StudentIDs: []uuid.UUID{studentID},
		Present:    true,
		Notes:      ptrString("Student participated actively in today's session"),
	}

	db := repo.GetDB()
	results, err := repo.CreateSessionStudent(ctx, db, input)
	assert.NoError(t, err)
	assert.NotNil(t, results)
	for idx, result := range *results {
		assert.Equal(t, idx+1, result.ID)
		assert.Equal(t, sessionID, result.SessionID)
		assert.Equal(t, studentID, result.StudentID)
		assert.True(t, result.Present)
		assert.NotNil(t, result.Notes)
		assert.Equal(t, "Student participated actively in today's session", *result.Notes)
		assert.False(t, result.CreatedAt.IsZero())
		assert.False(t, result.UpdatedAt.IsZero())
	}

	// Test duplicate creation (should fail due to unique constraint)
	duplicateInput := &models.CreateSessionStudentInput{
		SessionIDs: []uuid.UUID{sessionID},
		StudentIDs: []uuid.UUID{studentID},
		Present:    false,
		Notes:      ptrString("Duplicate entry"),
	}

	duplicateResult, err := repo.CreateSessionStudent(ctx, db, duplicateInput)
	assert.Error(t, err)
	assert.Nil(t, duplicateResult)

	// Test with invalid session ID (foreign key violation)
	invalidSessionID := uuid.New()
	invalidInput := &models.CreateSessionStudentInput{
		SessionIDs: []uuid.UUID{invalidSessionID},
		StudentIDs: []uuid.UUID{studentID},
		Present:    *ptrBool(false),
		Notes:      nil,
	}

	invalidResult, err := repo.CreateSessionStudent(ctx, db, invalidInput)
	assert.Error(t, err)
	assert.Nil(t, invalidResult)
}

func TestSessionStudentRepository_DeleteSessionStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionStudentRepository(testDB)
	ctx := context.Background()

	// Create test data using helper functions
	therapistID := CreateTestTherapist(t, testDB, ctx)
	sessionID := CreateTestSession(t, testDB, ctx, therapistID, "Session for delete test")
	studentID := CreateTestStudent(t, testDB, ctx, therapistID, "Delete")

	// Create session-student relationship
	_, err := testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id, present, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, sessionID, studentID, true, "Initial relationship")
	assert.NoError(t, err)

	// Test successful deletion
	deleteInput := &models.DeleteSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
	}

	err = repo.DeleteSessionStudent(ctx, deleteInput)
	assert.NoError(t, err)

	// Verify deletion - should not exist anymore
	var count int
	err = testDB.QueryRow(ctx, `
		SELECT COUNT(*) FROM session_student 
		WHERE session_id = $1 AND student_id = $2
	`, sessionID, studentID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestSessionStudentRepository_PatchSessionStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionStudentRepository(testDB)
	ctx := context.Background()

	// Create test data using helper functions
	therapistID := CreateTestTherapist(t, testDB, ctx)
	sessionID := CreateTestSession(t, testDB, ctx, therapistID, "Session for patch test")
	studentID := CreateTestStudent(t, testDB, ctx, therapistID, "Patch")

	// Create initial session-student relationship
	_, err := testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id, present, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, sessionID, studentID, true, "Original notes")
	assert.NoError(t, err)

	// Test patching present field only
	presentFalse := false
	patchInput := &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   &presentFalse,
		Notes:     nil,
	}

	result, err := repo.PatchSessionStudent(ctx, patchInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, studentID, result.StudentID)
	assert.False(t, result.Present)
	assert.NotNil(t, result.Notes)
	assert.Equal(t, "Original notes", *result.Notes)

	// Test patching notes field only
	newNotes := ptrString("Updated progress notes")
	patchInput = &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Notes:     newNotes,
	}

	result, err = repo.PatchSessionStudent(ctx, patchInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated progress notes", *result.Notes)

	// Test patching both fields
	presentTrue := true
	bothFieldsNotes := ptrString("Final comprehensive notes")
	patchInput = &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   &presentTrue,
		Notes:     bothFieldsNotes,
	}

	result, err = repo.PatchSessionStudent(ctx, patchInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Present)
	assert.Equal(t, "Final comprehensive notes", *result.Notes)

	// Test patching non-existent relationship
	nonExistentInput := &models.PatchSessionStudentInput{
		SessionID: uuid.New(),
		StudentID: uuid.New(),
		Present:   &presentTrue,
		Notes:     ptrString("Should fail"),
	}

	failResult, err := repo.PatchSessionStudent(ctx, nonExistentInput)
	assert.Error(t, err)
	assert.Nil(t, failResult)
}

func TestSessionStudentRepository_RateStudentSession(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionStudentRepository(testDB)
	ctx := context.Background()

	// Create test data using helper functions
	therapistID := CreateTestTherapist(t, testDB, ctx)
	sessionID := CreateTestSession(t, testDB, ctx, therapistID, "Session for rating test")
	studentID := CreateTestStudent(t, testDB, ctx, therapistID, "Rating")

	// Create initial session-student relationship
	_, err := testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id, present, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, sessionID, studentID, true, "Initial notes")
	require.NoError(t, err)

	// Test 1: Create new ratings
	presentTrue := true
	notes := ptrString("Session went well")
	rateInput := &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   &presentTrue,
		Notes:     notes,
		Ratings: &[]models.RateInput{
			{
				Category:    "visual_cue",
				Level:       "minimal",
				Description: "Minimal visual prompting needed",
			},
			{
				Category:    "engagement",
				Level:       "high",
				Description: "Highly engaged",
			},
		},
	}

	sessionStudent, ratings, err := repo.RateStudentSession(ctx, rateInput)
	require.NoError(t, err)
	require.NotNil(t, sessionStudent)
	assert.Equal(t, sessionID, sessionStudent.SessionID)
	assert.Equal(t, studentID, sessionStudent.StudentID)
	assert.True(t, sessionStudent.Present)
	assert.Equal(t, "Session went well", *sessionStudent.Notes)
	assert.Len(t, ratings, 2)

	// Test 2: Update existing ratings
	rateInput = &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   nil,
		Notes:     ptrString("Updated notes"),
		Ratings: &[]models.RateInput{
			{
				Category:    "visual_cue",
				Level:       "maximal",
				Description: "Required significant visual support",
			},
		},
	}

	sessionStudent, ratings, err = repo.RateStudentSession(ctx, rateInput)
	require.NoError(t, err)
	assert.Equal(t, "Updated notes", *sessionStudent.Notes)
	assert.Len(t, ratings, 1)
	assert.Equal(t, "visual_cue", *ratings[0].Category)
	assert.Equal(t, "maximal", *ratings[0].Level)

	// Test 3: Empty ratings array
	emptyRatingsInput := &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   &presentTrue,
		Notes:     ptrString("No ratings update"),
		Ratings:   &[]models.RateInput{},
	}

	sessionStudent, ratings, err = repo.RateStudentSession(ctx, emptyRatingsInput)
	require.NoError(t, err)
	assert.Equal(t, "No ratings update", *sessionStudent.Notes)
	assert.Len(t, ratings, 0)

	// Test 4: Non-existent session-student relationship
	nonExistentInput := &models.PatchSessionStudentInput{
		SessionID: uuid.New(),
		StudentID: uuid.New(),
		Present:   ptrBool(true),
		Notes:     ptrString("Should fail"),
		Ratings: &[]models.RateInput{
			{
				Category:    "visual_cue",
				Level:       "minimal",
				Description: "Should not be inserted",
			},
		},
	}

	failSessionStudent, failRatings, err := repo.RateStudentSession(ctx, nonExistentInput)
	assert.Error(t, err)
	assert.Nil(t, failSessionStudent)
	assert.Nil(t, failRatings)
}
