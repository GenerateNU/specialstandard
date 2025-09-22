package schema_test

import (
	"context"
	"testing"
	"time"

	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ptrInt(i int) *int {
	return &i
}

func TestThemeRepository_CreateTheme(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewThemeRepository(testDB.Pool)
	ctx := context.Background()

	// Test successful creation
	input := &models.CreateThemeInput{
		Name:  "Spring 2024",
		Month: 3,
		Year:  2024,
	}

	theme, err := repo.CreateTheme(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, theme)
	
	if theme != nil {
		assert.NotEqual(t, uuid.Nil, theme.ID)
		assert.Equal(t, "Spring 2024", theme.Name)
		assert.Equal(t, 3, theme.Month)
		assert.Equal(t, 2024, theme.Year)
		assert.NotNil(t, theme.CreatedAt)
		assert.NotNil(t, theme.UpdatedAt)
	}

	// Verify the theme was actually inserted into the database
	if theme != nil {
		var insertedTheme models.Theme
		err = testDB.Pool.QueryRow(ctx, `
			SELECT id, theme_name, month, year, created_at, updated_at 
			FROM theme WHERE id = $1
		`, theme.ID).Scan(
			&insertedTheme.ID,
			&insertedTheme.Name,
			&insertedTheme.Month,
			&insertedTheme.Year,
			&insertedTheme.CreatedAt,
			&insertedTheme.UpdatedAt,
		)
		assert.NoError(t, err)
		assert.Equal(t, theme.ID, insertedTheme.ID)
		assert.Equal(t, "Spring 2024", insertedTheme.Name)
	}
}

func TestThemeRepository_GetThemes(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewThemeRepository(testDB.Pool)
	ctx := context.Background()

	// Insert test themes
	theme1ID := uuid.New()
	theme2ID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6), ($7, $8, $9, $10, $11, $12)
    `, theme1ID, "Spring 2024", 3, 2024, time.Now(), time.Now(),
		theme2ID, "Summer 2024", 6, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	themes, err := repo.GetThemes(ctx)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, themes, 2)
	
	// Themes should be ordered by year DESC, month DESC
	// Both are 2024, so Summer (month 6) should come before Spring (month 3)
	assert.Equal(t, "Summer 2024", themes[0].Name)
	assert.Equal(t, "Spring 2024", themes[1].Name)
	assert.Equal(t, theme2ID, themes[0].ID)
	assert.Equal(t, theme1ID, themes[1].ID)

	// Test empty result
	_, err = testDB.Pool.Exec(ctx, "DELETE FROM theme")
	assert.NoError(t, err)
	
	themes, err = repo.GetThemes(ctx)
	assert.NoError(t, err)
	assert.Len(t, themes, 0)
}

func TestThemeRepository_GetThemeByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewThemeRepository(testDB.Pool)
	ctx := context.Background()

	// Insert test theme
	themeID := uuid.New()
	testTime := time.Now()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "Autumn 2024", 9, 2024, testTime, testTime)
	assert.NoError(t, err)

	// Test successful retrieval
	theme, err := repo.GetThemeByID(ctx, themeID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, theme)
	assert.Equal(t, themeID, theme.ID)
	assert.Equal(t, "Autumn 2024", theme.Name)
	assert.Equal(t, 9, theme.Month)
	assert.Equal(t, 2024, theme.Year)
	assert.NotNil(t, theme.CreatedAt)
	assert.NotNil(t, theme.UpdatedAt)

	// Test not found
	nonExistentID := uuid.New()
	theme, err = repo.GetThemeByID(ctx, nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, theme)
	assert.Contains(t, err.Error(), "Error querying database for given ID")
}

func TestThemeRepository_PatchTheme(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewThemeRepository(testDB.Pool)
	ctx := context.Background()

	// Insert test theme
	themeID := uuid.New()
	testTime := time.Now()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "Winter 2024", 12, 2024, testTime, testTime)
	assert.NoError(t, err)

	// Test successful patch with name only
	input := &models.UpdateThemeInput{
		Name: ptrString("Winter Wonderland 2024"),
	}

	theme, err := repo.PatchTheme(ctx, themeID, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, theme)
	assert.Equal(t, themeID, theme.ID)
	assert.Equal(t, "Winter Wonderland 2024", theme.Name)
	assert.Equal(t, 12, theme.Month) // Unchanged
	assert.Equal(t, 2024, theme.Year) // Unchanged

	// Test patch with multiple fields
	input2 := &models.UpdateThemeInput{
		Month: ptrInt(1),
		Year:  ptrInt(2025),
	}

	theme, err = repo.PatchTheme(ctx, themeID, input2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Winter Wonderland 2024", theme.Name) // Unchanged from previous update
	assert.Equal(t, 1, theme.Month) // Updated
	assert.Equal(t, 2025, theme.Year) // Updated

	// Test not found
	nonExistentID := uuid.New()
	input3 := &models.UpdateThemeInput{
		Name: ptrString("Non-existent"),
	}
	theme, err = repo.PatchTheme(ctx, nonExistentID, input3)
	assert.Error(t, err)
	assert.Nil(t, theme)
	assert.Contains(t, err.Error(), "error querying database for given theme ID")

	// Test no fields to update
	emptyInput := &models.UpdateThemeInput{}
	theme, err = repo.PatchTheme(ctx, themeID, emptyInput)
	assert.Error(t, err)
	assert.Nil(t, theme)
	assert.Contains(t, err.Error(), "No fields given to update")
}

func TestThemeRepository_DeleteTheme(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewThemeRepository(testDB.Pool)
	ctx := context.Background()

	// Insert test theme
	themeID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "To Be Deleted", 5, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test successful deletion
	err = repo.DeleteTheme(ctx, themeID)
	assert.NoError(t, err)

	// Verify theme was deleted
	var count int
	err = testDB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM theme WHERE id = $1`, themeID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// Test idempotent deletion (theme already deleted)
	err = repo.DeleteTheme(ctx, themeID)
	assert.NoError(t, err) // Should not error for non-existent theme

	// Test deletion of non-existent theme
	nonExistentID := uuid.New()
	err = repo.DeleteTheme(ctx, nonExistentID)
	assert.NoError(t, err) // Should not error for non-existent theme (idempotent)
}

func TestThemeRepository_DatabaseConstraints(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewThemeRepository(testDB.Pool)
	ctx := context.Background()

	// Test month constraint (should be between 1-12)
	input := &models.CreateThemeInput{
		Name:  "Invalid Month Theme",
		Month: 13, // Invalid month
		Year:  2024,
	}

	theme, err := repo.CreateTheme(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, theme)

	// Test year constraint (should be >= 1900)
	input2 := &models.CreateThemeInput{
		Name:  "Invalid Year Theme",
		Month: 6,
		Year:  1800, // Invalid year
	}

	theme, err = repo.CreateTheme(ctx, input2)
	assert.Error(t, err)
	assert.Nil(t, theme)

	// Test valid constraints
	input3 := &models.CreateThemeInput{
		Name:  "Valid Theme",
		Month: 6,
		Year:  2024,
	}

	theme, err = repo.CreateTheme(ctx, input3)
	assert.NoError(t, err)
	assert.NotNil(t, theme)
}
