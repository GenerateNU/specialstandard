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

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestStudentRepository_GetStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapists
	therapistID1 := uuid.New()
	therapistID2 := uuid.New()

	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID1, "Kevin", "Matula", "matulakevin91@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID2, "Jane", "Smith", "janesmith@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test students with different grades and therapists
	studentID1 := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2010-05-15")
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, studentID1, "John", "Doe", testDOB, therapistID1, 5, "IEP Goals: Speech articulation", time.Now(), time.Now())
	assert.NoError(t, err)

	studentID2 := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, studentID2, "Jane", "Smith", testDOB, therapistID2, 3, "IEP Goals: Reading", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test 1: Get all students (no filters)
	students, err := repo.GetStudents(ctx, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	// Test 2: Filter by grade
	students, err = repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Doe", students[0].LastName)
	assert.Equal(t, studentID1, students[0].ID)

	// Test 3: Filter by therapist
	students, err = repo.GetStudents(ctx, nil, therapistID2, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)
	assert.Equal(t, therapistID2, students[0].TherapistID)

	// Test 4: Filter by name (first name)
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)

	// Test 5: Filter by name (last name)
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "Smith", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)

	// Test 6: Multiple filters
	students, err = repo.GetStudents(ctx, ptrInt(5), therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)
	assert.Equal(t, 5, *students[0].Grade)

	// Test 7: Filter that returns no results
	students, err = repo.GetStudents(ctx, ptrInt(99), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 0)

	// More Tests for Pagination Behaviour
	for i := 3; i <= 6; i++ {
		testDOB, _ := time.Parse("2006-01-02", "2004-09-24")
		_, err := testDB.Pool.Exec(ctx, `
            INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            `, uuid.New(), "Student", fmt.Sprintf("Number%d", i), testDOB, therapistID1, i, "IEP: GOALS!!", time.Now(), time.Now())
		assert.NoError(t, err)
	}

	// Test 8: Pagination - get all students
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 6) // 2 original + 4 new

	// Test 9: Pagination - second page
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "", utils.Pagination{
		Page:  2,
		Limit: 5,
	})
	assert.NoError(t, err)
	assert.Len(t, students, 1)
}

// Add these additional test functions to your student_test.go file in the schema package

func TestStudentRepository_GetStudents_FilterByGrade(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Test", "Therapist", "test@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create students with different grades
	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "John", "Doe", testDOB, therapistID, 5, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Jane", "Smith", testDOB, therapistID, 4, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Mike", "Johnson", testDOB, therapistID, 5, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test: Filter by grade 5
	students, err := repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.NewPagination())

	assert.NoError(t, err)
	assert.Len(t, students, 2) // Should only return John and Mike
	for _, student := range students {
		assert.Equal(t, 5, *student.Grade)
	}
}

func TestStudentRepository_GetStudents_FilterByTherapist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create two therapists
	therapistID1 := uuid.New()
	therapistID2 := uuid.New()

	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID1, "Therapist", "One", "one@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID2, "Therapist", "Two", "two@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create students assigned to different therapists
	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Student", "One", testDOB, therapistID1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Student", "Two", testDOB, therapistID1, 4, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Student", "Three", testDOB, therapistID2, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test: Filter by therapist 1
	students, err := repo.GetStudents(ctx, nil, therapistID1, "", utils.NewPagination())

	assert.NoError(t, err)
	assert.Len(t, students, 2) // Should only return students assigned to therapist 1
	for _, student := range students {
		assert.Equal(t, therapistID1, student.TherapistID)
	}

	// Test: Filter by therapist 2
	students, err = repo.GetStudents(ctx, nil, therapistID2, "", utils.NewPagination())

	assert.NoError(t, err)
	assert.Len(t, students, 1) // Should only return student assigned to therapist 2
	assert.Equal(t, therapistID2, students[0].TherapistID)
}

func TestStudentRepository_GetStudents_FilterByName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Test", "Therapist", "test@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create students with different names
	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "John", "Doe", testDOB, therapistID, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Jane", "Johnson", testDOB, therapistID, 4, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Michael", "Johns", testDOB, therapistID, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test: Search by first name
	students, err := repo.GetStudents(ctx, nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // Should find "John" and "Johnson"

	// Test: Search by last name
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "Doe", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Doe", students[0].LastName)

	// Test: Case insensitive search
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "JOHN", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // Should find "John", "Johnson", and "Johns"

	// Test: Partial name search
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "oh", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // Should find "John", "Johnson", and "Johns"
}

func TestStudentRepository_GetStudents_CombinedFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create two therapists
	therapistID1 := uuid.New()
	therapistID2 := uuid.New()

	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID1, "Therapist", "One", "one@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID2, "Therapist", "Two", "two@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create diverse student data
	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	// Student 1: Therapist 1, Grade 5, Name John
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "John", "Doe", testDOB, therapistID1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	// Student 2: Therapist 1, Grade 4, Name Jane
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Jane", "Smith", testDOB, therapistID1, 4, time.Now(), time.Now())
	assert.NoError(t, err)

	// Student 3: Therapist 2, Grade 5, Name John
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "John", "Wilson", testDOB, therapistID2, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	// Student 4: Therapist 1, Grade 5, Name Sarah
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "Sarah", "Johnson", testDOB, therapistID1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test 1: Filter by grade only
	students, err := repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // John Doe, John Wilson, Sarah Johnson
	for _, student := range students {
		assert.Equal(t, 5, *student.Grade)
	}

	// Test 2: Filter by therapist only
	students, err = repo.GetStudents(ctx, nil, therapistID1, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // John Doe, Jane Smith, Sarah Johnson
	for _, student := range students {
		assert.Equal(t, therapistID1, student.TherapistID)
	}

	// Test 3: Filter by name only
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // John Doe, John Wilson, Sarah Johnson (has "John" in last name)

	// Test 4: Combine grade + therapist
	students, err = repo.GetStudents(ctx, ptrInt(5), therapistID1, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2) // John Doe, Sarah Johnson
	for _, student := range students {
		assert.Equal(t, 5, *student.Grade)
		assert.Equal(t, therapistID1, student.TherapistID)
	}

	// Test 5: Combine grade + name
	students, err = repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3) // John Doe, John Wilson, Sarah Johnson

	// Test 6: Combine therapist + name
	students, err = repo.GetStudents(ctx, nil, therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2) // John Doe, Sarah Johnson

	// Test 7: All filters combined
	students, err = repo.GetStudents(ctx, ptrInt(5), therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2) // John Doe, Sarah Johnson

	// Test 8: Filters that return no results
	students, err = repo.GetStudents(ctx, ptrInt(12), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 0)
}

func TestStudentRepository_GetStudents_CaseInsensitiveSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Test", "Therapist", "test@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test student
	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, uuid.New(), "John", "Smith", testDOB, therapistID, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	// Test lowercase search
	students, err := repo.GetStudents(ctx, nil, uuid.Nil, "john", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)

	// Test uppercase search
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "SMITH", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)

	// Test mixed case search
	students, err = repo.GetStudents(ctx, nil, uuid.Nil, "JoHn", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
}

func TestStudentRepository_GetStudents_WithPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Test", "Therapist", "test@test.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create 6 students all with grade 5
	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 6; i++ {
		_, err = testDB.Pool.Exec(ctx, `
            INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `, uuid.New(), fmt.Sprintf("Student%d", i), "Test", testDOB, therapistID, 5, time.Now(), time.Now())
		assert.NoError(t, err)
	}

	// Test: First page with limit 3
	students, err := repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.Pagination{Page: 1, Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	// Test: Second page with limit 3
	students, err = repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.Pagination{Page: 2, Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	// Test: Third page with limit 3 (should be empty or partial)
	students, err = repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.Pagination{Page: 3, Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, students, 0)

	// Test: Get all with large limit
	students, err = repo.GetStudents(ctx, ptrInt(5), uuid.Nil, "", utils.Pagination{Page: 1, Limit: 100})
	assert.NoError(t, err)
	assert.Len(t, students, 6)
}

func TestStudentRepository_GetStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist first
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2010-05-15")
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, studentID, "Jane", "Smith", testDOB, therapistID, 3, "IEP Goals: Language comprehension", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	student, err := repo.GetStudent(ctx, studentID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Smith", student.LastName)
	assert.Equal(t, "Jane", student.FirstName)
	assert.Equal(t, studentID, student.ID)
	assert.Equal(t, therapistID, student.TherapistID)
	assert.Equal(t, 3, *student.Grade)
}

func TestStudentRepository_AddStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist first
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test student data (no manual ID)
	testDOB, _ := time.Parse("2006-01-02", "2012-08-20")
	testStudent := models.Student{
		// ID removed - let database generate it
		FirstName:   "Alex",
		LastName:    "Johnson",
		DOB:         ptrTime(testDOB),
		TherapistID: therapistID,
		Grade:       ptrInt(2),
		IEP:         ptrString("IEP Goals: Fluency and articulation"),
	}

	// Test - get the database-generated student
	createdStudent, err := repo.AddStudent(ctx, testStudent)
	assert.NoError(t, err)

	// Verify student was inserted correctly using the returned ID
	var insertedStudent models.Student
	err = testDB.Pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, dob, therapist_id, grade, iep 
		FROM student WHERE id = $1
	`, createdStudent.ID).Scan( // Use the returned ID, not a manual one
		&insertedStudent.ID,
		&insertedStudent.FirstName,
		&insertedStudent.LastName,
		&insertedStudent.DOB,
		&insertedStudent.TherapistID,
		&insertedStudent.Grade,
		&insertedStudent.IEP,
	)
	assert.NoError(t, err)
	assert.Equal(t, testStudent.LastName, insertedStudent.LastName)
	assert.Equal(t, testStudent.TherapistID, insertedStudent.TherapistID)

	// Verify the ID was actually generated by the database
	assert.NotEqual(t, uuid.Nil, createdStudent.ID)
	assert.Equal(t, createdStudent.ID, insertedStudent.ID)
}

func TestStudentRepository_UpdateStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist first
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2011-03-10")
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, studentID, "Sam", "Wilson", testDOB, therapistID, 4, "Initial IEP", time.Now(), time.Now())
	assert.NoError(t, err)

	// Update student data
	updatedStudent := models.Student{
		ID:          studentID,
		FirstName:   "Sam",
		LastName:    "Wilson-Updated",
		DOB:         ptrTime(testDOB),
		TherapistID: therapistID,
		Grade:       ptrInt(5),
		IEP:         ptrString("Updated IEP Goals: Advanced speech therapy"),
	}

	// Test
	_, err = repo.UpdateStudent(ctx, updatedStudent)

	// Assert
	assert.NoError(t, err)

	// Verify student was updated correctly
	var verifyStudent models.Student
	err = testDB.Pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, dob, therapist_id, grade, iep 
		FROM student WHERE id = $1
	`, studentID).Scan(
		&verifyStudent.ID,
		&verifyStudent.FirstName,
		&verifyStudent.LastName,
		&verifyStudent.DOB,
		&verifyStudent.TherapistID,
		&verifyStudent.Grade,
		&verifyStudent.IEP,
	)
	assert.NoError(t, err)
	assert.Equal(t, "Wilson-Updated", verifyStudent.LastName)
	assert.Equal(t, 5, *verifyStudent.Grade)
}

func TestStudentRepository_DeleteStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	// Setup
	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	repo := schema.NewStudentRepository(testDB.Pool)
	ctx := context.Background()

	// Create test therapist first
	therapistID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, time.Now(), time.Now())
	assert.NoError(t, err)

	// Create test student
	studentID := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2009-12-25")
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, studentID, "Chris", "Brown", testDOB, therapistID, 6, "IEP Goals: Social communication", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	err = repo.DeleteStudent(ctx, studentID)

	// Assert
	assert.NoError(t, err)

	// Verify student was deleted
	var count int
	err = testDB.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM student WHERE id = $1`, studentID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestStudentRepository_PromoteStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Database Test in Short Mode")
	}

	testDB := testutil.SetupTestDB(t)
	defer testDB.Cleanup()

	ctx := context.Background()
	repo := schema.NewStudentRepository(testDB.Pool)

	// Insert 2 Therapists
	doctorWhoID := uuid.New()
	doctorDoofenshmirtzID := uuid.New()
	_, err := testDB.Pool.Exec(ctx, `
		INSERT INTO therapist (id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4),
		       ($5, $6, $7, $8);
    `, doctorWhoID, "Doctor", "Who", "doc.who@guesswho.com",
		doctorDoofenshmirtzID, "Heinz", "Doofenshmirtz", "doofy.balooney@gmail.com")
	assert.NoError(t, err)

	// Insert 5 Students
	_, err = testDB.Pool.Exec(ctx, `
		INSERT INTO student (first_name, last_name, dob, therapist_id, grade, iep, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW()),
		       ($7, $8, $3, $4, $9, $6, NOW()),
		       ($10, $11, $3, $4, $12, $6, NOW()),
		       ($13, $14, $3, $4, $15, $6, NOW()),
		       ($16, $17, $3, $18, $5, $6, NOW());
    `, "Michelle", "Li", "2005-01-31", doctorWhoID, 0, "IEP Content",
		"Ally", "Descoteaux", 5,
		"Luis", "Enrique Sarmiento", 12,
		"Harsh", "Singh", -1,
		"Stone", "Liu", doctorDoofenshmirtzID)
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, 5, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)

	patch := models.PromoteStudentsInput{
		TherapistID:        doctorWhoID,
		ExcludedStudentIDs: []uuid.UUID{students[0].ID},
	}
	err = repo.PromoteStudents(ctx, patch)
	assert.NoError(t, err)
	// TODO (Note for Future-Harsh): Might need to change this test call, after merging branch and comment correction.
	students, err = repo.GetStudents(ctx, 0, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)

	for i := 0; i < len(students); i++ {
		switch {
		case students[i].FirstName == "Michelle":
			assert.Equal(t, *students[i].Grade, 1)
		case students[i].FirstName == "Ally":
			assert.Equal(t, *students[i].Grade, 5)
		case students[i].FirstName == "Luis" || students[i].FirstName == "Harsh":
			assert.Equal(t, *students[i].Grade, -1)
		case students[i].FirstName == "Stone":
			assert.Equal(t, *students[i].Grade, 0)
		}
	}
}
