package schema_test

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGameContentRepository_GetGameContents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewGameContentRepository(testDB)
	ctx := context.Background()

	themeID := uuid.New()
	input := models.GetGameContentRequest{
		ThemeID:         ptrUUID(themeID),
		Category:        ptrString("speech"),
		QuestionType:    ptrString("sequencing"),
		DifficultyLevel: ptrInt(5),
		QuestionCount:   ptrInt(4),
		WordsCount:      ptrInt(3),
	}

	// Empty DB
	gameContents, err := repo.GetGameContents(ctx, input)
	assert.Empty(t, gameContents)
	assert.NoError(t, err)

	id := uuid.New()
	category := "speech"
	week := 3
	questionType := "sequencing"
	level := 5
	question := "Go watch Bugonia? It's the best movie EVER!"
	options := []string{"Koorkodile", "Krokorok", "Korkorockodile"}
	answer := "Crocodile"

	// Inserting GameContent (FK error with theme_id)
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, theme_id, week, category, question_type, difficulty_level, 
		                          question, options, answer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, id, themeID, 3, category, questionType, level, question, options, answer)
	assert.Error(t, err)

	// Inserting Valid ThemeID.
	_, err = testDB.Exec(ctx, `
		INSERT INTO theme (id, theme_name, month, year)
		VALUES ($1, $2, $3, $4)
    `, themeID, "Animal Kingdom", 4, 2019)
	assert.NoError(t, err)

	// Inserting GameContent Finally!
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, theme_id, week, category, question_type, difficulty_level, 
		                          question, options, answer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, id, themeID, 3, category, questionType, level, question, options, answer)
	assert.NoError(t, err)

	gameContents, err = repo.GetGameContents(ctx, input)
	assert.NotNil(t, gameContents)
	assert.NoError(t, err)

	gameContent := gameContents[0]

	assert.Equal(t, gameContent.ID, id)
	assert.Equal(t, gameContent.ThemeID, themeID)
	assert.Equal(t, gameContent.Week, week)
	assert.Equal(t, gameContent.Category, &category)
	assert.Equal(t, gameContent.QuestionType, questionType)
	assert.Equal(t, gameContent.DifficultyLevel, level)
	assert.Equal(t, gameContent.Question, question)
	for _, word := range gameContent.Options {
		assert.Contains(t, gameContent.Options, word)
	}
	assert.Equal(t, len(gameContent.Options), len(options)-1)
	assert.Equal(t, gameContent.Answer, answer)
}
