package schema_test

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestGameContentRepository_GetGameContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewGameContentRepository(testDB)
	ctx := context.Background()

	input := models.GetGameContentRequest{
		Category:        "sequencing",
		DifficultyLevel: 5,
		Count:           4,
	}

	// Empty DB
	gameContent, err := repo.GetGameContent(ctx, input)
	assert.Nil(t, gameContent)
	assert.Error(t, err)
	assert.ErrorIs(t, err, pgx.ErrNoRows)

	id := uuid.New()
	category := "sequencing"
	level := 5
	options := []string{"Koorkodile", "Krokorok", "Krookodile", "Korkorockodile"}
	answer := "Crocodile"

	// Inserting GameContent!
	_, err = testDB.Exec(ctx, `
		INSERT INTO game_content (id, category, level, options, answer)
		VALUES ($1, $2, $3, $4, $5)
    `, id, category, level, options, answer)
	assert.NoError(t, err)

	gameContent, err = repo.GetGameContent(ctx, input)
	assert.NotNil(t, gameContent)
	assert.NoError(t, err)

	assert.Equal(t, gameContent.ID, id)
	assert.Equal(t, gameContent.Category, category)
	assert.Equal(t, gameContent.DifficultyLevel, level)
	for _, word := range gameContent.Options {
		assert.Contains(t, gameContent.Options, word)
	}
	assert.Equal(t, len(gameContent.Options), len(options)-1)
	assert.Equal(t, gameContent.Category, category)
}
