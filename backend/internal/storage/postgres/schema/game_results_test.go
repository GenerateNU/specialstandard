package schema_test

import (
	"context"
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

	// Inserting Valid Therapist
	therapistID := uuid.New()
	_, err := testDB.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email)
        	 VALUES ($1, $2, $3, $4)`,
		therapistID, "Speech", "Therapist", "teachthespeech@specialstandard.com")
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Test session for session-student")
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, therapist_id, grade)
        VALUES ($1, $2, $3, $4, $5)
    `, studentID, "Test", "Student", therapistID, 5)
	assert.NoError(t, err)

	// Inserting SessionStudent
	_, err = testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id)
		VALUES ($1, $2);
    `, sessionID, studentID)
	assert.NoError(t, err)

	// Inserting GameContent!
	contentID := uuid.New()
	category := "sequencing"
	level := 5
	options := []string{"Koorkodile", "Krokorok", "Krookodile", "Korkorockodile"}
	answer := "Crocodile"
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, category, level, options, answer)
		VALUES ($1, $2, $3, $4, $5)
    `, contentID, category, level, options, answer)
	assert.NoError(t, err)

	// Inserting GameResult
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_result (session_id, student_id, content_id, time_taken, completed, incorrect_tries)
		VALUES ($1, $2, $3, $4, $5, $6);
    `, sessionID, studentID, contentID, 93, true, 5)
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

	assert.Equal(t, gameResult.SessionID, sessionID)
	assert.Equal(t, gameResult.StudentID, studentID)
	assert.Equal(t, gameResult.ContentID, contentID)
	assert.Equal(t, gameResult.TimeTaken, 93)
	assert.True(t, gameResult.Completed)
	assert.Equal(t, gameResult.IncorrectTries, 5)
}

func TestGameResultRepository_PostGameResult(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewGameResultRepository(testDB)
	ctx := context.Background()

	input := models.PostGameResult{
		SessionID:      uuid.New(),
		StudentID:      uuid.New(),
		ContentID:      uuid.New(),
		TimeTakenSec:   90,
		Completed:      ptr.Bool(true),
		IncorrectTries: ptr.Int(5),
	}

	postedGameResult, err := repo.PostGameResult(ctx, input)
	assert.Nil(t, postedGameResult)
	assert.Error(t, err)

	// Inserting Valid Therapist
	therapistID := uuid.New()
	_, err = testDB.Exec(ctx,
		`INSERT INTO therapist (id, first_name, last_name, email)
        	 VALUES ($1, $2, $3, $4)`,
		therapistID, "Speech", "Therapist", "teachthespeech@specialstandard.com")
	assert.NoError(t, err)

	// Create test session
	sessionID := uuid.New()
	startTime := time.Now()
	endTime := startTime.Add(time.Hour)
	_, err = testDB.Exec(ctx, `
        INSERT INTO session (id, therapist_id, start_datetime, end_datetime, notes)
        VALUES ($1, $2, $3, $4, $5)
    `, sessionID, therapistID, startTime, endTime, "Test session for session-student")
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, therapist_id, grade)
        VALUES ($1, $2, $3, $4, $5)
    `, studentID, "Test", "Student", therapistID, 5)
	assert.NoError(t, err)

	// Inserting SessionStudent
	_, err = testDB.Exec(ctx, `
		INSERT INTO session_student (session_id, student_id)
		VALUES ($1, $2);
    `, sessionID, studentID)
	assert.NoError(t, err)

	// Inserting GameContent!
	contentID := uuid.New()
	category := "sequencing"
	level := 5
	options := []string{"Koorkodile", "Krokorok", "Krookodile", "Korkorockodile"}
	answer := "Crocodile"
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, category, level, options, answer)
		VALUES ($1, $2, $3, $4, $5)
    `, contentID, category, level, options, answer)
	assert.NoError(t, err)

	input = models.PostGameResult{
		SessionID:      sessionID,
		StudentID:      studentID,
		ContentID:      contentID,
		TimeTakenSec:   -6,
		Completed:      ptr.Bool(true),
		IncorrectTries: ptr.Int(4),
	}

	postedGameResult, err = repo.PostGameResult(ctx, input)
	assert.Nil(t, postedGameResult)
	assert.Error(t, err)

	input.TimeTakenSec = 91

	postedGameResult, err = repo.PostGameResult(ctx, input)
	assert.NotNil(t, postedGameResult)
	assert.NoError(t, err)

	assert.Equal(t, postedGameResult.SessionID, sessionID)
	assert.Equal(t, postedGameResult.StudentID, studentID)
	assert.Equal(t, postedGameResult.ContentID, contentID)
	assert.Equal(t, postedGameResult.TimeTaken, 91)
	assert.Equal(t, postedGameResult.Completed, true)
	assert.Equal(t, postedGameResult.IncorrectTries, 4)
}
