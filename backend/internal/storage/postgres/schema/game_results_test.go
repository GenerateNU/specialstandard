package schema_test

import (
	"context"
	"math/rand"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"
	"specialstandard/internal/utils"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ptrUUID(id uuid.UUID) *uuid.UUID {
	return &id
}

func TestGameResultRepository_GetGameResults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewGameResultRepository(testDB)
	ctx := context.Background()

	schoolID := 1
	_, err := testDB.Exec(ctx,
		`INSERT INTO school (id, name)
	   	 VALUES ($1, $2)`,
		schoolID, "Special Standard School")
	assert.NoError(t, err)

	// Inserting Valid Therapist
	therapistID := uuid.New()
	_, err = testDB.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email, schools)
       	 VALUES ($1, $2, $3, $4, $5)`,
		therapistID, "Speech", "Therapist", "teachthespeech@specialstandard.com", []int{schoolID})
	assert.NoError(t, err)

	// Create test session
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

	sessionID := uuid.New()
	_, err = testDB.Exec(ctx, `
	   INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, session_parent_id)
	   VALUES ($1, $2, $3, $4, $5, $6)
   `, sessionID, "Test Session", startTime, endTime, "Test session for session-student", sessionParentID)
	assert.NoError(t, err)
	// Convert datetimes to dates (zero out time portion) and update inserted rows,

	_, err = testDB.Exec(ctx, `
		UPDATE session_parent
		SET start_date = $1, end_date = $2
		WHERE id = $3
	`, startDate, endDate, sessionParentID)
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
		UPDATE session
		SET start_datetime = $1, end_datetime = $2
		WHERE id = $3
	`, startDate, endDate, sessionID)
	assert.NoError(t, err)
	// Create test student
	studentID := uuid.New()
	_, err = testDB.Exec(ctx, `
       INSERT INTO student (id, first_name, last_name, therapist_id, school_id, grade)
       VALUES ($1, $2, $3, $4, $5, $6)
   `, studentID, "Test", "Student", therapistID, schoolID, 5)
	assert.NoError(t, err)

	// Inserting SessionStudent
	sessionStudentID := rand.Intn(10)
	_, err = testDB.Exec(ctx, `
		INSERT INTO session_student (id, session_id, student_id)
		VALUES ($1, $2, $3);
   `, sessionStudentID, sessionID, studentID)
	assert.NoError(t, err)

	// Inserting valid Theme
	themeID := uuid.New()
	_, err = testDB.Exec(ctx, `
		INSERT INTO theme (id, theme_name, month, year)
		VALUES ($1, $2, $3, $4)
    `, themeID, "Animalia", 4, 2018)
	assert.NoError(t, err)

	// Inserting GameContent
	contentID := uuid.New()
	category := "speech"
	week := 3
	questionType := "sequencing"
	level := 5
	question := "Go watch Bugonia? It's the best movie EVER!"
	options := []string{"Koorkodile", "Krokorok", "Korkorockodile"}
	answer := "Crocodile"
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, theme_id, week, category, question_type, difficulty_level, 
		                          question, options, answer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, contentID, themeID, week, category, questionType, level, question, options, answer)
	assert.NoError(t, err)

	// Inserting GameResult
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_result (session_student_id, content_id, time_taken_sec, completed, 
		                         count_of_incorrect_attempts, incorrect_attempts)
		VALUES ($1, $2, $3, $4, $5, $6);
   `, sessionStudentID, contentID, 93, true, 5, []string{"Kroorockodile"})
	assert.NoError(t, err)

	gameResultQuery := &models.GetGameResultQuery{
		SessionID: ptrUUID(sessionID),
		StudentID: ptrUUID(studentID),
	}
	gameResults, err := repo.GetGameResults(ctx, gameResultQuery, utils.NewPagination())

	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.NotNil(t, gameResults)
	assert.Equal(t, len(gameResults), 1)

	gameResult := gameResults[0]

	assert.Equal(t, gameResult.SessionStudentID, sessionStudentID)
	assert.Equal(t, gameResult.ContentID, contentID)
	assert.Equal(t, gameResult.TimeTakenSec, 93)
	assert.True(t, gameResult.Completed)
	assert.Equal(t, gameResult.IncorrectAttempts, &[]string{"Kroorockodile"})
	assert.Equal(t, gameResult.CountIncorrectAttempts, 5)
}

func TestGameResultRepository_PostGameResult(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewGameResultRepository(testDB)
	ctx := context.Background()

	input := models.PostGameResult{
		SessionStudentID:       rand.Intn(10),
		ContentID:              uuid.New(),
		TimeTakenSec:           90,
		Completed:              ptr.Bool(true),
		CountIncorrectAttempts: 5,
	}

	postedGameResult, err := repo.PostGameResult(ctx, input)
	assert.Nil(t, postedGameResult)
	assert.Error(t, err)

	schoolID := 1
	_, err = testDB.Exec(ctx,
		`INSERT INTO school (id, name)
	   	 VALUES ($1, $2)`,
		schoolID, "Special Standard School")
	assert.NoError(t, err)

	// Inserting Valid Therapist
	therapistID := uuid.New()
	_, err = testDB.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email, schools)
       	 VALUES ($1, $2, $3, $4, $5)`,
		therapistID, "Speech", "Therapist", "teachthespeech@specialstandard.com", []int{schoolID})
	assert.NoError(t, err)

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

	// Create test session
	sessionID := uuid.New()
	_, err = testDB.Exec(ctx, `
       INSERT INTO session (id, session_name, start_datetime, end_datetime, notes, session_parent_id)
       VALUES ($1, $2, $3, $4, $5, $6)
   `, sessionID, "Test Session", startTime, endTime, "Test session for session-student", sessionParentID)
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Exec(ctx, `
       INSERT INTO student (id, first_name, last_name, therapist_id, school_id, grade)
       VALUES ($1, $2, $3, $4, $5, $6)
   `, studentID, "Test", "Student", therapistID, schoolID, 5)
	assert.NoError(t, err)

	// Inserting SessionStudent
	sessionStudentID := rand.Intn(10)
	_, err = testDB.Exec(ctx, `
		INSERT INTO session_student (id, session_id, student_id)
		VALUES ($1, $2, $3);
   `, sessionStudentID, sessionID, studentID)
	assert.NoError(t, err)

	// Inserting Theme
	themeID := uuid.New()
	_, err = testDB.Exec(ctx, `
		INSERT INTO theme (id, theme_name, month, year)
		VALUES ($1, $2, $3, $4)
    `, themeID, "Animalia", 4, 2018)
	assert.NoError(t, err)

	// Inserting GameContent!
	contentID := uuid.New()
	category := "speech"
	level := 5
	options := []string{"Koorkodile", "Krokorok", "Krookodile", "Korkorockodile"}
	answer := "Crocodile"
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, theme_id, week, category, question_type, difficulty_level, 
		                          question, options, answer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
   `, contentID, themeID, 2, category, "sequencing", level, "Are you shtruggling with your esshesh?", options, answer)
	assert.NoError(t, err)

	input = models.PostGameResult{
		SessionStudentID:       sessionStudentID,
		ContentID:              contentID,
		TimeTakenSec:           -6,
		Completed:              ptr.Bool(true),
		CountIncorrectAttempts: 4,
	}

	postedGameResult, err = repo.PostGameResult(ctx, input)
	assert.Nil(t, postedGameResult)
	assert.Error(t, err)

	input.TimeTakenSec = 91

	postedGameResult, err = repo.PostGameResult(ctx, input)
	assert.NotNil(t, postedGameResult)
	assert.NoError(t, err)

	assert.Equal(t, postedGameResult.SessionStudentID, sessionStudentID)
	assert.Equal(t, postedGameResult.ContentID, contentID)
	assert.Equal(t, postedGameResult.TimeTakenSec, 91)
	assert.Equal(t, postedGameResult.Completed, true)
	assert.Equal(t, postedGameResult.CountIncorrectAttempts, 4)
	assert.Equal(t, postedGameResult.IncorrectAttempts, &[]string{})
}
