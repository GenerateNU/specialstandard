package schema_test

import (
	"context"
	"fmt"
	"specialstandard/internal/utils"
	"testing"
	"time"

	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/storage/postgres/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestResourceRepository_GetResources(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewResourceRepository(testDB.Pool)
	ctx := context.Background()

	// Create test theme first (foreign key requirement)
	themeID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "Test Theme", 1, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test resource
	resourceID := uuid.New()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, resourceID, themeID, "5th Grade", testDate, "worksheet", "Math Worksheet", "mathematics", "Addition problems", time.Now(), time.Now())
	assert.NoError(t, err)

	resources, err := repo.GetResources(ctx, themeID, "", "", "", "", "", "", nil, 0, 0, utils.NewPagination())

	// Assert
	assert.NoError(t, err)
	print(len(resources))
	if assert.NotEmpty(t, resources, "Expected at least 1 resource") {
		assert.Equal(t, "Math Worksheet", *resources[0].Title)
		assert.Equal(t, resourceID, resources[0].ID)
		assert.Equal(t, themeID, resources[0].ThemeID)
		assert.Equal(t, "Test Theme", resources[0].Theme.Name)
		assert.Equal(t, 1, resources[0].Theme.Month)
		assert.Equal(t, 2024, resources[0].Theme.Year)
	}

	// More Tests for Pagination Behaviour
	for i := 2; i <= 15; i++ {
		_, err := testDB.Pool.Exec(ctx,
			`INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, uuid.New(), themeID, fmt.Sprintf("%d-th Grade", i), testDate, "worksheet", "Math Worksheet", "mathematics", "Addition problems", time.Now(), time.Now())
		assert.NoError(t, err)
	}

	resources, err = repo.GetResources(ctx, themeID, "", "", "", "", "", "", nil, 0, 0, utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, resources, 10)

	resources, err = repo.GetResources(ctx, themeID, "", "", "", "", "", "", nil, 0, 0, utils.Pagination{
		Page:  2,
		Limit: 9,
	})
	assert.NoError(t, err)
	assert.Len(t, resources, 6)
	assert.Equal(t, "10-th Grade", *resources[0].GradeLevel)
}

func TestResourceRepository_GetResourceByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewResourceRepository(testDB.Pool)
	ctx := context.Background()

	// Create test theme first
	themeID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "Science Theme", 2, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test resource
	resourceID := uuid.New()
	testDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, resourceID, themeID, "3rd Grade", testDate, "video", "Science Video", "science", "Volcano experiment", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	resource, err := repo.GetResourceByID(ctx, resourceID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resource)
	if resource != nil {
		assert.Equal(t, "Science Video", *resource.Title)
		assert.Equal(t, resourceID, resource.ID)
		assert.Equal(t, themeID, resource.ThemeID)
		assert.Equal(t, "Science Theme", resource.Theme.Name)
	}
}

func TestResourceRepository_CreateResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewResourceRepository(testDB.Pool)
	ctx := context.Background()

	// Create test theme first
	themeID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "Reading Theme", 3, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	testDate := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	resourceBody := models.ResourceBody{
		ThemeID:    themeID,
		GradeLevel: ptrString("4th Grade"),
		Date:       &testDate,
		Type:       ptrString("book"),
		Title:      ptrString("Reading Comprehension"),
		Category:   ptrString("literacy"),
		Content:    ptrString("Story analysis exercises"),
	}

	// Test
	resource, err := repo.CreateResource(ctx, resourceBody)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resource)
	if resource != nil {
		assert.Equal(t, "Reading Comprehension", *resource.Title)
		assert.Equal(t, themeID, resource.ThemeID)
	}
}

func TestResourceRepository_UpdateResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewResourceRepository(testDB.Pool)
	ctx := context.Background()

	// Create test theme first
	themeID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "Art Theme", 4, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test resource
	resourceID := uuid.New()
	testDate := time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, resourceID, themeID, "2nd Grade", testDate, "activity", "Art Activity", "art", "Drawing exercise", time.Now(), time.Now())
	assert.NoError(t, err)

	updateBody := models.UpdateResourceBody{
		Title:    ptrString("Updated Art Activity"),
		Category: ptrString("creative"),
	}

	// Test
	resource, err := repo.UpdateResource(ctx, resourceID, updateBody)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resource)
	if resource != nil {
		assert.Equal(t, "Updated Art Activity", *resource.Title)
		assert.Equal(t, resourceID, resource.ID)
	}
}

func TestResourceRepository_DeleteResource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewResourceRepository(testDB.Pool)
	ctx := context.Background()

	// Create test theme first
	themeID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO theme (id, theme_name, month, year, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, themeID, "History Theme", 5, 2024, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test resource
	resourceID := uuid.New()
	testDate := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO resource (id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, resourceID, themeID, "6th Grade", testDate, "document", "History Document", "history", "Ancient civilizations", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	err = repo.DeleteResource(ctx, resourceID)

	// Assert
	assert.NoError(t, err)

	var count int
	err = testDB.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM resource WHERE id = $1", resourceID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
