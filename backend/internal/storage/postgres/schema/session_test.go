package schema_test

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"testing"
	"time"

	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ptrString(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

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
	sessions, err := repo.GetSessions(ctx, utils.NewPagination(), nil)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, "Test session", *sessions[0].Notes)
	assert.Equal(t, therapistID, sessions[0].TherapistID)
	assert.True(t, sessions[0].EndDateTime.After(sessions[0].StartDateTime))

	// More Tests for Pagination Behaviour
	for i := 1; i <= 18; i++ {
		start := startTime.Add(time.Duration(i) * time.Hour)
		end := start.Add(time.Hour)

		_, err := testDB.Pool.Exec(ctx, `
			INSERT INTO session (therapist_id, start_datetime, end_datetime, notes)
			VALUES ($1, $2, $3, $4)
       `, therapistID, start, end, fmt.Sprintf("Test session%d", i))
		assert.NoError(t, err)
	}

	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), nil)

	assert.NoError(t, err)
	assert.Len(t, sessions, 10)

	sessions, err = repo.GetSessions(ctx, utils.Pagination{
		Page:  4,
		Limit: 5,
	}, nil)

	assert.NoError(t, err)
	assert.Len(t, sessions, 4)
	assert.Equal(t, "Test session", *sessions[3].Notes)

	// Test filtering by year
	yearFilter := &models.GetSessionRepositoryRequest{
		Year: intPtr(startTime.Year()),
	}
	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), yearFilter)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(sessions))

	// Test filtering by month and year
	monthYearFilter := &models.GetSessionRepositoryRequest{
		Month: intPtr(int(startTime.Month())),
		Year:  intPtr(startTime.Year()),
	}
	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), monthYearFilter)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(sessions))

	// Test filtering by student IDs 
	studentID1 := uuid.New()
	studentID2 := uuid.New()
	
	// Insert student associations for one of the sessions
	sessionWithStudents := sessions[0].ID
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO student (id, first_name, last_name, therapist_id)
		VALUES ($1, $2, $3, $4), ($5, $6, $7, $8)
	`, studentID1, "Student", "One", therapistID, studentID2, "Student", "Two", therapistID)
	assert.NoError(t, err)
	
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id, present)
		VALUES ($1, $2, true), ($3, $4, true)
	`, sessionWithStudents, studentID1, sessionWithStudents, studentID2)
	assert.NoError(t, err)

	studentFilter := &models.GetSessionRepositoryRequest{
		StudentIDs: &[]uuid.UUID{studentID1, studentID2},
	}
	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), studentFilter)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
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

	// SUCCESS TEST - Creation of valid Therapist first.
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx,
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

	err = repo.DeleteSession(ctx, sessionID)
	assert.NoError(t, err)
}

func TestSessionRepository_PostSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB tests in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionRepository(testDB.Pool)
	ctx := context.Background()

	therapistID := uuid.New()
	startTime := time.Now()
	endTime := time.Now().Add(time.Hour)
	notes := ptrString("foreign key violation")
	postSession := &models.PostSessionInput{
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       notes,
	}
	postedSession, err := repo.PostSession(ctx, postSession)
	assert.Error(t, err)
	assert.Nil(t, postedSession)

	// INSERTING VALID THERAPIST
	therapistID = uuid.New()
	_, err = testDB.Pool.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email)
        	 VALUES ($1, $2, $3, $4)`,
		therapistID, "Speech", "Therapist", "teachthespeech@specialstandard.com")
	assert.NoError(t, err)

	startTime = time.Now()
	endTime = time.Now().Add(-time.Hour)
	notes = ptrString("check constraint violation")
	postSession = &models.PostSessionInput{
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       notes,
	}
	postedSession, err = repo.PostSession(ctx, postSession)
	assert.Error(t, err)
	assert.Nil(t, postedSession)
	assert.False(t, endTime.After(startTime))

	startTime = time.Now()
	endTime = time.Now().Add(time.Hour)
	notes = ptrString("success")
	postSession = &models.PostSessionInput{
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       notes,
	}
	postedSession, err = repo.PostSession(ctx, postSession)
	assert.NoError(t, err)
	assert.NotNil(t, postedSession)
	assert.Equal(t, postedSession.TherapistID, therapistID)
	assert.Equal(t, postedSession.Notes, notes)
	assert.True(t, postedSession.EndDateTime.After(postedSession.StartDateTime))
}

func TestSessionRepository_PatchSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB Tests in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewSessionRepository(testDB.Pool)
	ctx := context.Background()

	// Given ID Not Found 404 Error
	badID := uuid.New()
	patch := &models.PatchSessionInput{
		Notes: ptrString("404 NOT FOUND ERROR"),
	}
	patchedSession, err := repo.PatchSession(ctx, badID, patch)
	assert.Error(t, err)
	assert.Nil(t, patchedSession)

	// Foreign Key Violation
	id := uuid.New()
	therapistID := uuid.New()
	patch = &models.PatchSessionInput{
		TherapistID: &therapistID,
	}
	patchedSession, err = repo.PatchSession(ctx, id, patch)
	assert.Error(t, err)
	assert.Nil(t, patchedSession)

	// INSERTING THERAPIST NOW
	therapistID = uuid.New()
	_, err = testDB.Pool.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email)
             VALUES ($1, $2, $3, $4)`,
		therapistID, "Doc", "The Dwarf", "doc@sevendwarves.com")
	assert.NoError(t, err)

	startTime := time.Now()
	endTime := time.Now().Add(-time.Hour)
	notes := ptrString("check constraint violation")
	patch = &models.PatchSessionInput{
		StartTime: &startTime,
		EndTime:   &endTime,
		Notes:     notes,
	}
	patchedSession, err = repo.PatchSession(ctx, id, patch)
	assert.Error(t, err)
	assert.Nil(t, patchedSession)
	assert.False(t, endTime.After(startTime))

	// INSERT ACTUAL SESSION TO EDIT
	id = uuid.New()
	startTime = time.Now()
	endTime = time.Now().Add(time.Hour)
	_, err = testDB.Pool.Exec(ctx,
		`INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
             VALUES ($1, $2, $3, $4, $5)`,
		id, therapistID, startTime, endTime, "Inserted")
	assert.NoError(t, err)

	notes = ptrString("success with one change")
	patch = &models.PatchSessionInput{
		Notes: notes,
	}
	patchedSession, err = repo.PatchSession(ctx, id, patch)
	assert.NoError(t, err)
	assert.NotNil(t, patchedSession)
	assert.True(t, patchedSession.EndDateTime.After(patchedSession.StartDateTime))
	assert.Equal(t, patchedSession.TherapistID, therapistID)
	assert.Equal(t, patchedSession.Notes, notes)

	startTime = time.Now()
	endTime = time.Now().Add(time.Hour)
	patch = &models.PatchSessionInput{
		StartTime: &startTime,
		EndTime:   &endTime,
	}
	patchedSession, err = repo.PatchSession(ctx, id, patch)
	assert.NoError(t, err)
	assert.NotNil(t, patchedSession)
	assert.True(t, patchedSession.EndDateTime.After(patchedSession.StartDateTime))
	assert.Equal(t, patchedSession.TherapistID, therapistID)
	assert.Equal(t, patchedSession.Notes, notes)

	// ADDING A SECOND THERAPIST TO UPDATE TO
	therapistID = uuid.New()
	_, err = testDB.Pool.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email)
             VALUES ($1, $2, $3, $4)`,
		therapistID, "Courage", "The Cowardly Dog", "havecourage@cowardice.com")
	assert.NoError(t, err)

	startTime = time.Now()
	endTime = time.Now().Add(time.Hour)
	notes = ptrString("New Note")
	patch = &models.PatchSessionInput{
		StartTime:   &startTime,
		EndTime:     &endTime,
		TherapistID: &therapistID,
		Notes:       notes,
	}
	patchedSession, err = repo.PatchSession(ctx, id, patch)
	assert.NoError(t, err)
	assert.NotNil(t, patchedSession)
	assert.True(t, patchedSession.EndDateTime.After(patchedSession.StartDateTime))
	assert.Equal(t, patchedSession.TherapistID, therapistID)
	assert.Equal(t, patchedSession.Notes, notes)
}
