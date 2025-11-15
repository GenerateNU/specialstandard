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

// CreateSessionTestTherapist creates a therapist with all required fields for session tests
func CreateSessionTestTherapist(t *testing.T, db *pgxpool.Pool, ctx context.Context, name string) uuid.UUID {
	// Ensure district and school exist
	districtID := 1
	_, err := db.Exec(ctx, `
		INSERT INTO district (id, name, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, districtID, "Test District")
	assert.NoError(t, err)

	schoolID := 1
	_, err = db.Exec(ctx, `
		INSERT INTO school (id, name, district_id, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, schoolID, "Test School", districtID)
	assert.NoError(t, err)

	therapistID := uuid.New()
	email := fmt.Sprintf("%s_%s@example.com", name, therapistID.String()[:8])
	_, err = db.Exec(ctx, `
		INSERT INTO therapist (id, first_name, last_name, email, schools, district_id, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`, therapistID, name, "Therapist", email, []int{schoolID}, districtID, true)
	assert.NoError(t, err)

	return therapistID
}

// CreateSessionTestStudent creates a student for session tests
func CreateSessionTestStudent(t *testing.T, db *pgxpool.Pool, ctx context.Context, therapistID uuid.UUID, name string, grade int) uuid.UUID {
	// Ensure school exists for student
	districtID := 1
	_, err := db.Exec(ctx, `
		INSERT INTO district (id, name, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, districtID, "Test District")
	assert.NoError(t, err)

	schoolID := 1
	_, err = db.Exec(ctx, `
		INSERT INTO school (id, name, district_id, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, schoolID, "Test School", districtID)
	assert.NoError(t, err)

	studentID := uuid.New()
	_, err = db.Exec(ctx, `
		INSERT INTO student (id, first_name, last_name, therapist_id, school_id, grade, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`, studentID, name, "Student", therapistID, schoolID, grade)
	assert.NoError(t, err)

	return studentID
}

func TestSessionRepository_GetSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionRepository(testDB)
	ctx := context.Background()

	// Create a test therapist using helper
	therapistID := CreateSessionTestTherapist(t, testDB, ctx, "John")

	// Insert test session data using new schema
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
        INSERT INTO session (session_name, start_datetime, end_datetime, notes, location, created_at, updated_at, session_parent_id)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6)
    `, "Session Name", startTime, endTime, "Test session", "Centre of the Earth", sessionParentID)
	assert.NoError(t, err)

	// Test
	sessions, err := repo.GetSessions(ctx, utils.NewPagination(), nil, therapistID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)
	assert.Equal(t, "Test session", *sessions[0].Notes)
	assert.True(t, sessions[0].EndDateTime.After(sessions[0].StartDateTime))
	assert.Equal(t, "Session Name", sessions[0].SessionName)
	assert.Equal(t, "Centre of the Earth", *sessions[0].Location)

	// More Tests for Pagination Behaviour
	for i := 1; i <= 18; i++ {
		start := startTime.Add(time.Duration(i) * time.Hour)
		end := start.Add(time.Hour)

		sessionParentID := uuid.New()
		startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
		endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
		_, err := testDB.Exec(ctx, `
       INSERT INTO session_parent (id, start_date, end_date, therapist_id)
       VALUES ($1, $2, $3, $4)
   `, sessionParentID, startDate, endDate, therapistID)
		assert.NoError(t, err)

		_, err = testDB.Exec(ctx, `
			INSERT INTO session (session_name, start_datetime, end_datetime, notes, location, created_at, updated_at, session_parent_id)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6)
       `, "Session Name", start, end, fmt.Sprintf("Test session%d", i), "Khol De Baahein ~ Monali Thakur", sessionParentID)
		assert.NoError(t, err)
	}

	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), nil, therapistID)

	assert.NoError(t, err)
	assert.Len(t, sessions, 19)

	sessions, err = repo.GetSessions(ctx, utils.Pagination{
		Page:  4,
		Limit: 5,
	}, nil, therapistID)

	assert.NoError(t, err)
	assert.Len(t, sessions, 4)
	assert.Equal(t, "Test session18", *sessions[3].Notes)

	// Test filtering by year
	yearFilter := &models.GetSessionRepositoryRequest{
		Year: ptrInt(startTime.Year()),
	}
	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), yearFilter, therapistID)
	assert.NoError(t, err)
	assert.Equal(t, 19, len(sessions))

	// Test filtering by month and year
	monthYearFilter := &models.GetSessionRepositoryRequest{
		Month: ptrInt(int(startTime.Month())),
		Year:  ptrInt(startTime.Year()),
	}
	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), monthYearFilter, therapistID)
	assert.NoError(t, err)
	assert.Equal(t, 19, len(sessions))

	// Test filtering by student IDs
	studentID1 := CreateSessionTestStudent(t, testDB, ctx, therapistID, "Student1", 5)
	studentID2 := CreateSessionTestStudent(t, testDB, ctx, therapistID, "Student2", 5)

	// Insert student associations for one of the sessions
	sessionWithStudents := sessions[0].ID
	_, err = testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id, present, created_at, updated_at)
		VALUES ($1, $2, true, NOW(), NOW()), ($3, $4, true, NOW(), NOW())
	`, sessionWithStudents, studentID1, sessionWithStudents, studentID2)
	assert.NoError(t, err)

	studentFilter := &models.GetSessionRepositoryRequest{
		StudentIDs: &[]uuid.UUID{studentID1, studentID2},
	}
	sessions, err = repo.GetSessions(ctx, utils.NewPagination(), studentFilter, therapistID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(sessions))
}

func TestSessionRepository_GetSessionByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionRepository(testDB)
	ctx := context.Background()

	// Create a test therapist using helper
	therapistID := CreateSessionTestTherapist(t, testDB, ctx, "Jane")

	// Insert test session and capture the generated ID
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
        INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, location, created_at, updated_at, session_parent_id)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7)
    `, sessionID, "Session Name", startTime, endTime, "Get by ID test session", "Featurethon!", sessionParentID)
	assert.NoError(t, err)

	// Test
	session, err := repo.GetSessionByID(ctx, sessionID.String())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, sessionID, session.ID)
	//assert.Equal(t, therapistID, session.TherapistID)
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

	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionRepository(testDB)
	ctx := context.Background()

	// Create valid therapist using helper
	therapistID := CreateSessionTestTherapist(t, testDB, ctx, "Doctor")

	// Inserting test session into the DB
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

	_, err = testDB.Exec(ctx,
		`INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, location, created_at, updated_at, session_parent_id)
             VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7)`,
		sessionID, "Session Name", startTime, endTime, "Inserting into session for test", "Calcutta", sessionParentID)
	assert.NoError(t, err)

	err = repo.DeleteSession(ctx, sessionID)
	assert.NoError(t, err)
}

func TestSessionRepository_PostSessions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping DB tests in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionRepository(testDB)
	ctx := context.Background()

	// Test foreign key violation with non-existent therapist
	therapistID := uuid.New()
	sessionName := "Session Name"
	startTime := time.Now()
	endTime := time.Now().Add(time.Hour)
	notes := ptrString("foreign key violation")
	location := ptrString("Home is where the heart is, but God I love the English")
	postSession := &models.PostSessionInput{
		SessionName: sessionName,
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       notes,
		Location:    location,
	}
	db := repo.GetDB()
	postedSession, err := repo.PostSession(ctx, db, postSession)
	assert.Error(t, err)
	assert.Nil(t, postedSession)

	// Create valid therapist using helper
	therapistID = CreateSessionTestTherapist(t, testDB, ctx, "Speech")

	// Test check constraint violation
	startTime = time.Now()
	endTime = time.Now().Add(-time.Hour)
	notes = ptrString("check constraint violation")
	postSession = &models.PostSessionInput{
		SessionName: sessionName,
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       notes,
		Location:    location,
	}
	postedSession, err = repo.PostSession(ctx, db, postSession)
	assert.Error(t, err)
	assert.Nil(t, postedSession)
	assert.False(t, endTime.After(startTime))

	// Test successful session creation
	startTime = time.Now()
	endTime = time.Now().Add(time.Hour)
	notes = ptrString("success")
	postSession = &models.PostSessionInput{
		SessionName: sessionName,
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       notes,
		Location:    location,
	}
	postedSessions, err := repo.PostSession(ctx, db, postSession)
	assert.NoError(t, err)
	assert.NotNil(t, postedSessions)
	for _, postedSession := range *postedSessions {
		//assert.Equal(t, postedSession.TherapistID, therapistID)
		assert.Equal(t, postedSession.Notes, notes)
		assert.True(t, postedSession.EndDateTime.After(postedSession.StartDateTime))
	}

	// Test recurring sessions
	recurEnd := startTime.AddDate(0, 0, 20) // 3 weeks later
	postSession = &models.PostSessionInput{
		SessionName: sessionName,
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       ptrString("recurring sessions"),
		Location:    location,
		Repetition: &models.Repetition{
			EveryNWeeks: 1,
			RecurEnd:    recurEnd,
		},
	}

	repeatedSessions, err := repo.PostSession(ctx, db, postSession)
	assert.NoError(t, err)
	assert.NotNil(t, repeatedSessions)
	//assert.Equal(t, len(*repeatedSessions), 3)

	for _, s := range *repeatedSessions {
		//assert.Equal(t, s.TherapistID, therapistID)
		assert.Contains(t, *s.Notes, "recurring")
	}

	// Test invalid repetition
	postSession = &models.PostSessionInput{
		SessionName: sessionName,
		StartTime:   startTime,
		EndTime:     endTime,
		TherapistID: therapistID,
		Notes:       ptrString("invalid repetition end"),
		Repetition: &models.Repetition{
			EveryNWeeks: 1,
			RecurStart:  startTime,
			RecurEnd:    startTime.AddDate(0, 0, -7), // 1 week before start
		},
	}

	invalidRepeatSessions, err := repo.PostSession(ctx, db, postSession)
	assert.Error(t, err)
	assert.Nil(t, invalidRepeatSessions)
}

// func TestSessionRepository_PatchSessions(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping DB Tests in short mode")
// 	}

// 	testDB := testutil.SetupTestWithCleanup(t)
// 	repo := schema.NewSessionRepository(testDB)
// 	ctx := context.Background()

// 	// Test 404 not found
// 	badID := uuid.New()
// 	patch := &models.PatchSessionInput{
// 		Notes: ptrString("404 NOT FOUND ERROR"),
// 	}
// 	patchedSession, err := repo.PatchSession(ctx, badID, patch)
// 	assert.Error(t, err)
// 	assert.Nil(t, patchedSession)

// 	// Test foreign key violation
// 	id := uuid.New()
// 	therapistID := uuid.New()
// 	patch = &models.PatchSessionInput{
// 		TherapistID: &therapistID,
// 	}
// 	patchedSession, err = repo.PatchSession(ctx, id, patch)
// 	assert.Error(t, err)
// 	assert.Nil(t, patchedSession)

// 	// Create first therapist using helper
// 	therapistID = CreateSessionTestTherapist(t, testDB, ctx, "Doc")

// 	// Test check constraint violation
// 	startTime := time.Now()
// 	endTime := time.Now().Add(-time.Hour)
// 	notes := ptrString("check constraint violation")
// 	patch = &models.PatchSessionInput{
// 		StartTime: &startTime,
// 		EndTime:   &endTime,
// 		Notes:     notes,
// 	}
// 	patchedSession, err = repo.PatchSession(ctx, id, patch)
// 	assert.Error(t, err)
// 	assert.Nil(t, patchedSession)
// 	assert.False(t, endTime.After(startTime))

// 	// Insert actual session to edit
// 	id = uuid.New()
// 	startTime = time.Now()
// 	endTime = time.Now().Add(time.Hour)

// 	sessionParentID := uuid.New()
// 	startDate := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
// 	endDate := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 0, 0, 0, 0, endTime.Location())
// 	_, err = testDB.Exec(ctx, `
//        INSERT INTO session_parent (id, start_date, end_date, therapist_id)
//        VALUES ($1, $2, $3, $4)
//    `, sessionParentID, startDate, endDate, therapistID)
// 	assert.NoError(t, err)

// 	sessionName := "Test Session"
// 	location := "Area 51"
// 	_, err = testDB.Exec(ctx,
// 		`INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, location, created_at, updated_at, session_parent_id)
//              VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7)`,
// 		id, sessionName, startTime, endTime, "Inserted", location, sessionParentID)
// 	assert.NoError(t, err)

// 	// Test successful patch with notes only
// 	notes = ptrString("success with one change")
// 	patch = &models.PatchSessionInput{
// 		Notes: notes,
// 	}
// 	patchedSession, err = repo.PatchSession(ctx, id, patch)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, patchedSession)
// 	assert.True(t, patchedSession.EndDateTime.After(patchedSession.StartDateTime))
// 	assert.Equal(t, patchedSession.TherapistID, therapistID)
// 	assert.Equal(t, patchedSession.Notes, notes)

// 	// Test patch with time update
// 	startTime = time.Now()
// 	endTime = time.Now().Add(time.Hour)
// 	patch = &models.PatchSessionInput{
// 		StartTime: &startTime,
// 		EndTime:   &endTime,
// 	}
// 	patchedSession, err = repo.PatchSession(ctx, id, patch)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, patchedSession)
// 	assert.True(t, patchedSession.EndDateTime.After(patchedSession.StartDateTime))
// 	assert.Equal(t, patchedSession.TherapistID, therapistID)
// 	assert.Equal(t, patchedSession.Notes, notes)

// 	// Create second therapist using helper
// 	therapistID2 := CreateSessionTestTherapist(t, testDB, ctx, "Courage")

// 	// Test updating all fields
// 	startTime = time.Now()
// 	endTime = time.Now().Add(time.Hour)
// 	notes = ptrString("New Note")
// 	patch = &models.PatchSessionInput{
// 		SessionName: ptrString("New Test Session"),
// 		StartTime:   &startTime,
// 		EndTime:     &endTime,
// 		TherapistID: &therapistID2,
// 		Notes:       notes,
// 		Location:    ptrString("Area 52"),
// 	}
// 	patchedSession, err = repo.PatchSession(ctx, id, patch)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, patchedSession)
// 	assert.True(t, patchedSession.EndDateTime.After(patchedSession.StartDateTime))
// 	assert.Equal(t, patchedSession.TherapistID, therapistID2)
// 	assert.Equal(t, patchedSession.Notes, notes)
// 	assert.Equal(t, patchedSession.SessionName, "New Test Session")
// 	assert.Equal(t, *patchedSession.Location, "Area 52")
// }

func TestGetSessionStudents(t *testing.T) {
	// Setup
	testDB := testutil.SetupTestWithCleanup(t)
	repo := schema.NewSessionRepository(testDB)
	ctx := context.Background()

	// Create therapist using helper
	therapistID := CreateSessionTestTherapist(t, testDB, ctx, "John")

	// Insert test session
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
        INSERT INTO session (session_name, start_datetime, end_datetime, notes, created_at, updated_at, session_parent_id)
        VALUES ($1, $2, $3, $4, NOW(), NOW(), $5)
    `, "Test Session", startTime, endTime, "Test session", sessionParentID)
	assert.NoError(t, err)

	sessions, err := repo.GetSessions(ctx, utils.NewPagination(), nil, therapistID)
	assert.NoError(t, err)
	assert.Len(t, sessions, 1)

	// Create students using helper
	studentID1 := CreateSessionTestStudent(t, testDB, ctx, therapistID, "Student1", 3)
	studentID3 := CreateSessionTestStudent(t, testDB, ctx, therapistID, "Student3", -1) // graduated

	// Insert student associations
	sessionWithStudents := sessions[0].ID
	_, err = testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id, present, created_at, updated_at)
		VALUES ($1, $2, true, NOW(), NOW()), ($3, $4, true, NOW(), NOW())
	`, sessionWithStudents, studentID1, sessionWithStudents, studentID3)
	assert.NoError(t, err)

	students, err := repo.GetSessionStudents(ctx, sessionWithStudents, utils.NewPagination(), therapistID)

	assert.NoError(t, err)
	assert.Len(t, students, 1) // returns the one student that has not graduated
}
