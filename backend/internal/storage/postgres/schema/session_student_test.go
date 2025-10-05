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
)

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
		SessionID: sessionID,
		StudentID: studentID,
		Present:   true,
		Notes:     testutil.Ptr("Student participated actively in today's session"),
	}

	result, err := repo.CreateSessionStudent(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionID, result.SessionID)
	assert.Equal(t, studentID, result.StudentID)
	assert.True(t, result.Present)
	assert.NotNil(t, result.Notes)
	assert.Equal(t, "Student participated actively in today's session", *result.Notes)
	assert.False(t, result.CreatedAt.IsZero())
	assert.False(t, result.UpdatedAt.IsZero())

	// Test duplicate creation (should fail due to unique constraint)
	duplicateInput := &models.CreateSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   false,
		Notes:     testutil.Ptr("Duplicate entry"),
	}

	duplicateResult, err := repo.CreateSessionStudent(ctx, duplicateInput)
	assert.Error(t, err)
	assert.Nil(t, duplicateResult)

	// Test with invalid session ID (foreign key violation)
	invalidSessionID := uuid.New()
	invalidInput := &models.CreateSessionStudentInput{
		SessionID: invalidSessionID,
		StudentID: studentID,
		Present:   true,
		Notes:     nil,
	}

	invalidResult, err := repo.CreateSessionStudent(ctx, invalidInput)
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
	patchInput := &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   &presentFalse,
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
	newNotes := testutil.Ptr("Updated progress notes")
	patchInput = &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   nil, // Don't update present
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
	bothFieldsNotes := testutil.Ptr("Final comprehensive notes")
	patchInput = &models.PatchSessionStudentInput{
		SessionID: sessionID,
		StudentID: studentID,
		Present:   &presentTrue,
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
	nonExistentInput := &models.PatchSessionStudentInput{
		SessionID: uuid.New(),
		StudentID: uuid.New(),
		Present:   &presentTrue,
		Notes:     testutil.Ptr("Should fail"),
	}

	failResult, err := repo.PatchSessionStudent(ctx, nonExistentInput)
	assert.Error(t, err)
	assert.Nil(t, failResult)
}
