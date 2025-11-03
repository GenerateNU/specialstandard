package schema_test

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrBool(b bool) *bool {
	return &b
}

func TestSessionStudentRepository_CreateSessionStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist (required for session foreign key)
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email)
        VALUES ($1, $2, $3, $4)
    `, therapistID, "Dr. Test", "Therapist", "test.therapist@example.com")
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Test session for session-student")
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, therapist_id, grade)
        VALUES ($1, $2, $3, $4, $5)
    `, studentID, "Test", "Student", therapistID, 5)
	assert.NoError(t, err)

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
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email)
        VALUES ($1, $2, $3, $4)
    `, therapistID, "Dr. Delete", "Test", "delete.test@example.com")
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Session for delete test")
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, therapist_id, grade)
        VALUES ($1, $2, $3, $4, $5)
    `, studentID, "Delete", "Student", therapistID, 3)
	assert.NoError(t, err)

	// Create session-student relationship
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session_student (session_id, student_id, present, notes)
        VALUES ($1, $2, $3, $4)
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
	err = testDB.Pool.QueryRow(ctx, `
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
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email)
        VALUES ($1, $2, $3, $4)
    `, therapistID, "Dr. Patch", "Test", "patch.test@example.com")
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Session for patch test")
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, therapist_id, grade)
        VALUES ($1, $2, $3, $4, $5)
    `, studentID, "Patch", "Student", therapistID, 4)
	assert.NoError(t, err)

	// Create initial session-student relationship
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session_student (session_id, student_id, present, notes)
        VALUES ($1, $2, $3, $4)
    `, sessionID, studentID, true, "Original notes")
	assert.NoError(t, err)

	// Test patching present field only
	presentFalse := false
	patchInput := &models.SessionStudent{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   presentFalse,
		Notes:     nil, // Don't update notes
	}

	result, err := repo.PatchSessionStudent(ctx, patchInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, studentID, result.StudentID)
	assert.False(t, result.Present) // Should be updated
	assert.NotNil(t, result.Notes)
	assert.Equal(t, "Original notes", *result.Notes) // Should remain unchanged

	// Test patching notes field only
	newNotes := ptrString("Updated progress notes")
	patchInput = &models.SessionStudent{
		SessionID: sessionID,
		StudentID: studentID,
		Notes:     newNotes,
	}

	result, err = repo.PatchSessionStudent(ctx, patchInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, studentID, result.StudentID)
	assert.False(t, result.Present) // Should remain from previous update
	assert.NotNil(t, result.Notes)
	assert.Equal(t, "Updated progress notes", *result.Notes) // Should be updated

	// Test patching both fields
	presentTrue := true
	bothFieldsNotes := ptrString("Final comprehensive notes")
	patchInput = &models.SessionStudent{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   presentTrue,
		Notes:     bothFieldsNotes,
	}

	result, err = repo.PatchSessionStudent(ctx, patchInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, studentID, result.StudentID)
	assert.True(t, result.Present) // Should be updated
	assert.NotNil(t, result.Notes)
	assert.Equal(t, "Final comprehensive notes", *result.Notes) // Should be updated

	// Test patching non-existent relationship
	nonExistentInput := &models.SessionStudent{
		SessionID: uuid.New(),
		StudentID: uuid.New(),
		Present:   presentTrue,
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
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email)
        VALUES ($1, $2, $3, $4)
    `, therapistID, "Dr. Rate", "Test", "rate.test@example.com")
	require.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Session for rating test")
	require.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, therapist_id, grade)
        VALUES ($1, $2, $3, $4, $5)
    `, studentID, "Rating", "Student", therapistID, 5)
	require.NoError(t, err)

	// Create initial session-student relationship
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO session_student (session_id, student_id, present, notes)
        VALUES ($1, $2, $3, $4)
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

	// Test 2: Update existing ratings (ON CONFLICT)
	rateInput = &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   nil,
		Notes:     ptrString("Updated notes"),
		Ratings: &[]models.RateInput{
			{
				Category:    "visual_cue",
				Level:       "maximal", // Changed from minimal
				Description: "Required significant visual support",
			},
		},
	}

	sessionStudent, ratings, err = repo.RateStudentSession(ctx, rateInput)
	require.NoError(t, err)
	require.NotNil(t, sessionStudent)
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
	require.NotNil(t, sessionStudent)
	assert.Equal(t, "No ratings update", *sessionStudent.Notes)
	assert.NotNil(t, ratings)
	assert.Len(t, ratings, 0)

	// Test 4: Test engagement levels (low/high)
	engagementInput := &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   nil,
		Notes:     nil,
		Ratings: &[]models.RateInput{
			{
				Category:    "engagement",
				Level:       "low", // Testing low engagement
				Description: "Low engagement today",
			},
		},
	}

	sessionStudent, ratings, err = repo.RateStudentSession(ctx, engagementInput)
	require.NoError(t, err)
	assert.Len(t, ratings, 1)
	assert.NotNil(t, sessionStudent)
	assert.Equal(t, "engagement", *ratings[0].Category)
	assert.Equal(t, "low", *ratings[0].Level)

	// Test 5: Non-existent session-student relationship
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
