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

func PtrInt(i int) *int { return &i }

func TestStudentRepository_GetStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Generate Public Schools")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW()), ($4, $5, $6, NOW(), NOW());
	`, 1, "Generate Elementary", 1, 2, "Sherman Center for Academic Excellence", 1)
	assert.NoError(t, err)

	therapistID1 := uuid.New()
	therapistID2 := uuid.New()

	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID1, "Kevin", "Matula", "matulakevin91@gmail.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID2, "Jane", "Smith", "janesmith@gmail.com", true, []int{1, 2}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	studentID1 := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2010-05-15")
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, studentID1, "John", "Doe", testDOB, therapistID1, 1, 5, "IEP Goals: Speech articulation", time.Now(), time.Now())
	assert.NoError(t, err)

	studentID2 := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, studentID2, "Jane", "Smith", testDOB, therapistID2, 2, 3, "IEP Goals: Reading", time.Now(), time.Now())
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, nil, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Doe", students[0].LastName)
	assert.Equal(t, studentID1, students[0].ID)

	students, err = repo.GetStudents(ctx, nil, nil, therapistID2, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)
	assert.Equal(t, therapistID2, students[0].TherapistID)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "Smith", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)
	assert.Equal(t, 5, *students[0].Grade)

	students, err = repo.GetStudents(ctx, PtrInt(99), nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 0)

	for i := 3; i <= 6; i++ {
		testDOB, _ := time.Parse("2006-01-02", "2004-09-24")
		_, err := testDB.Exec(ctx, `
            INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            `, uuid.New(), "Student", fmt.Sprintf("Number%d", i), testDOB, therapistID1, 1, i, "IEP: GOALS!!", time.Now(), time.Now())
		assert.NoError(t, err)
	}

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 6)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "", utils.Pagination{
		Page:  2,
		Limit: 5,
	})
	assert.NoError(t, err)
	assert.Len(t, students, 1)

	students, err = repo.GetStudents(ctx, nil, PtrInt(1), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 5)

	students, err = repo.GetStudents(ctx, nil, PtrInt(2), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
}

func TestStudentRepository_GetStudents_FilterByGrade(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Test", "Therapist", "test@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, uuid.New(), "John", "Doe", testDOB, therapistID, 1, 5, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, uuid.New(), "Jane", "Smith", testDOB, therapistID, 1, 4, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, uuid.New(), "Mike", "Johnson", testDOB, therapistID, 1, 5, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, ` 
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, uuid.New(), "Jack", "Douglas", testDOB, therapistID, 1, -1, "IEP Goals", time.Now(), time.Now())
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, nil, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.NewPagination())

	assert.NoError(t, err)
	assert.Len(t, students, 2)
	for _, student := range students {
		assert.Equal(t, 5, *student.Grade)
	}

	students, err = repo.GetStudents(ctx, PtrInt(-1), nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, -1, *students[0].Grade)

}

func TestStudentRepository_GetStudents_FilterByTherapist(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID1 := uuid.New()
	therapistID2 := uuid.New()

	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID1, "Therapist", "One", "one@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID2, "Therapist", "Two", "two@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Student", "One", testDOB, therapistID1, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Student", "Two", testDOB, therapistID1, 1, 4, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Student", "Three", testDOB, therapistID2, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, nil, nil, therapistID1, "", utils.NewPagination())

	assert.NoError(t, err)
	assert.Len(t, students, 2)
	for _, student := range students {
		assert.Equal(t, therapistID1, student.TherapistID)
	}

	students, err = repo.GetStudents(ctx, nil, nil, therapistID2, "", utils.NewPagination())

	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, therapistID2, students[0].TherapistID)
}

func TestStudentRepository_GetStudents_FilterByName(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Test", "Therapist", "test@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "John", "Doe", testDOB, therapistID, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Jane", "Johnson", testDOB, therapistID, 1, 4, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Michael", "Johns", testDOB, therapistID, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, nil, nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "Doe", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Doe", students[0].LastName)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "JOHN", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "oh", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)
}

func TestStudentRepository_GetStudents_CombinedFilters(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID1 := uuid.New()
	therapistID2 := uuid.New()

	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID1, "Therapist", "One", "one@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID2, "Therapist", "Two", "two@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "John", "Doe", testDOB, therapistID1, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Jane", "Smith", testDOB, therapistID1, 1, 4, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "John", "Wilson", testDOB, therapistID2, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "Sarah", "Johnson", testDOB, therapistID1, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)
	for _, student := range students {
		assert.Equal(t, 5, *student.Grade)
	}

	students, err = repo.GetStudents(ctx, nil, nil, therapistID1, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)
	for _, student := range students {
		assert.Equal(t, therapistID1, student.TherapistID)
	}

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, therapistID1, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)
	for _, student := range students {
		assert.Equal(t, 5, *student.Grade)
		assert.Equal(t, therapistID1, student.TherapistID)
	}

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, nil, nil, therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	students, err = repo.GetStudents(ctx, PtrInt(12), nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 0)
}

func TestStudentRepository_GetStudents_CaseInsensitiveSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Test", "Therapist", "test@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, uuid.New(), "John", "Smith", testDOB, therapistID, 1, 5, time.Now(), time.Now())
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, nil, nil, uuid.Nil, "john", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "SMITH", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "JoHn", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
}

func TestStudentRepository_GetStudents_WithPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Test", "Therapist", "test@test.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 6; i++ {
		_, err = testDB.Exec(ctx, `
            INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `, uuid.New(), fmt.Sprintf("Student%d", i), "Test", testDOB, therapistID, 1, 5, time.Now(), time.Now())
		assert.NoError(t, err)
	}

	students, err := repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.Pagination{Page: 1, Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.Pagination{Page: 2, Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, students, 3)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.Pagination{Page: 3, Limit: 3})
	assert.NoError(t, err)
	assert.Len(t, students, 0)

	students, err = repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.Pagination{Page: 1, Limit: 100})
	assert.NoError(t, err)
	assert.Len(t, students, 6)
}

func TestStudentRepository_GetStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	studentID := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2010-05-15")
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, studentID, "Jane", "Smith", testDOB, therapistID, 1, 3, "IEP Goals: Language comprehension", time.Now(), time.Now())
	assert.NoError(t, err)

	student, err := repo.GetStudent(ctx, studentID)

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

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB, _ := time.Parse("2006-01-02", "2012-08-20")
	testStudent := models.Student{
		FirstName:   "Alex",
		LastName:    "Johnson",
		DOB:         ptrTime(testDOB),
		TherapistID: therapistID,
		SchoolID:    1,
		Grade:       PtrInt(2),
		IEP:         ptrString("IEP Goals: Fluency and articulation"),
	}

	createdStudent, err := repo.AddStudent(ctx, testStudent)
	assert.NoError(t, err)

	var insertedStudent models.Student
	err = testDB.QueryRow(ctx, `
		SELECT id, first_name, last_name, dob, therapist_id, school_id, grade, iep 
		FROM student WHERE id = $1
	`, createdStudent.ID).Scan(
		&insertedStudent.ID,
		&insertedStudent.FirstName,
		&insertedStudent.LastName,
		&insertedStudent.DOB,
		&insertedStudent.TherapistID,
		&insertedStudent.SchoolID,
		&insertedStudent.Grade,
		&insertedStudent.IEP,
	)
	assert.NoError(t, err)
	assert.Equal(t, testStudent.LastName, insertedStudent.LastName)
	assert.Equal(t, testStudent.TherapistID, insertedStudent.TherapistID)

	assert.NotEqual(t, uuid.Nil, createdStudent.ID)
	assert.Equal(t, createdStudent.ID, insertedStudent.ID)
}

func TestStudentRepository_UpdateStudent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	studentID := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2011-03-10")
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, studentID, "Sam", "Wilson", testDOB, therapistID, 1, 4, "Initial IEP", time.Now(), time.Now())
	assert.NoError(t, err)

	updatedStudent := models.Student{
		ID:          studentID,
		FirstName:   "Sam",
		LastName:    "Wilson-Updated",
		DOB:         ptrTime(testDOB),
		TherapistID: therapistID,
		SchoolID:    1,
		Grade:       PtrInt(5),
		IEP:         ptrString("Updated IEP Goals: Advanced speech therapy"),
	}

	_, err = repo.UpdateStudent(ctx, updatedStudent)

	assert.NoError(t, err)

	var verifyStudent models.Student
	err = testDB.QueryRow(ctx, `
		SELECT id, first_name, last_name, dob, therapist_id, school_id, grade, iep 
		FROM student WHERE id = $1
	`, studentID).Scan(
		&verifyStudent.ID,
		&verifyStudent.FirstName,
		&verifyStudent.LastName,
		&verifyStudent.DOB,
		&verifyStudent.TherapistID,
		&verifyStudent.SchoolID,
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

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Kevin", "Matula", "matulakevin91@gmail.com", true, []int{1}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	studentID := uuid.New()
	testDOB, _ := time.Parse("2006-01-02", "2009-12-25")
	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `, studentID, "Chris", "Brown", testDOB, therapistID, 1, 6, "IEP Goals: Social communication", time.Now(), time.Now())
	assert.NoError(t, err)

	err = repo.DeleteStudent(ctx, studentID)

	assert.NoError(t, err)

	var count int
	err = testDB.QueryRow(ctx, `SELECT COUNT(*) FROM student WHERE id = $1`, studentID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestStudentRepository_PromoteStudents(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Database Test in Short Mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	ctx := context.Background()
	repo := schema.NewStudentRepository(testDB)

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Test District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW());
	`, 1, "Test School", 1)
	assert.NoError(t, err)

	doctorWhoID := uuid.New()
	doctorDoofenshmirtzID := uuid.New()
	_, err = testDB.Exec(ctx, `
		INSERT INTO therapist (id, first_name, last_name, email, schools, district_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()),
		       ($7, $8, $9, $10, $11, $12, NOW(), NOW());
    `, doctorWhoID, "Doctor", "Who", "doc.who@guesswho.com", []int{1}, 1,
		doctorDoofenshmirtzID, "Heinz", "Doofenshmirtz", "doofy.balooney@gmail.com", []int{1}, 1)
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
		INSERT INTO student (first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()),
		       ($8, $9, $3, $4, $5, $10, $7, NOW()),
		       ($11, $12, $3, $4, $5, $13, $7, NOW()),
		       ($14, $15, $3, $4, $5, $16, $7, NOW()),
		       ($17, $18, $3, $19, $5, $6, $7, NOW());
    `, "Michelle", "Li", "2005-01-31", doctorWhoID, 1, 0, "IEP Content",
		"Ally", "Descoteaux", 5,
		"Luis", "Enrique Sarmiento", 12,
		"Harsh", "Singh", -1,
		"Stone", "Liu", doctorDoofenshmirtzID)
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, PtrInt(5), nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)

	patch := models.PromoteStudentsInput{
		TherapistID:        doctorWhoID,
		ExcludedStudentIDs: []uuid.UUID{students[0].ID},
	}
	err = repo.PromoteStudents(ctx, patch)
	assert.NoError(t, err)
	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)

	for i := 0; i < len(students); i++ {
		switch students[i].FirstName {
		case "Michelle":
			assert.Equal(t, *students[i].Grade, 1)
		case "Ally":
			assert.Equal(t, *students[i].Grade, 5)
		case "Luis":
			fallthrough
		case "Harsh":
			assert.Equal(t, *students[i].Grade, -1)
		case "Stone":
			assert.Equal(t, *students[i].Grade, 0)
		}
	}
}

func TestStudentRepository_GetStudents_FilterBySchoolAndDistrict(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW()), ($3, $4, NOW(), NOW());
	`, 1, "District One", 2, "District Two")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW()), ($4, $5, $6, NOW(), NOW()), ($7, $8, $9, NOW(), NOW());
	`, 1, "School A", 1, 2, "School B", 1, 3, "School C", 2)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "Test", "Therapist", "test@test.com", true, []int{1, 2, 3}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)

	_, err = testDB.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
        VALUES 
        ($1, $2, $3, $4, $5, $6, $7, $8, $9),
        ($10, $11, $12, $4, $5, $13, $7, $8, $9),
        ($14, $15, $16, $4, $5, $17, $7, $8, $9)
    `, uuid.New(), "Student", "One", testDOB, therapistID, 1, 5, time.Now(), time.Now(),
		uuid.New(), "Student", "Two", 2,
		uuid.New(), "Student", "Three", 3)
	assert.NoError(t, err)

	students, err := repo.GetStudents(ctx, nil, PtrInt(1), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "One", students[0].LastName)

	students, err = repo.GetStudents(ctx, nil, PtrInt(2), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Two", students[0].LastName)

	students, err = repo.GetStudents(ctx, nil, PtrInt(3), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Three", students[0].LastName)
}

func TestStudentRepository_GetStudents_MultipleSchoolsSameDistrict(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	testDB := testutil.SetupTestWithCleanup(t)

	repo := schema.NewStudentRepository(testDB)
	ctx := context.Background()

	_, err := testDB.Exec(ctx, `
	INSERT INTO "district" (id, name, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW());
	`, 1, "Unified District")
	assert.NoError(t, err)

	_, err = testDB.Exec(ctx, `
	INSERT INTO "school" (id, name, district_id, created_at, updated_at) 
	VALUES ($1, $2, $3, NOW(), NOW()), ($4, $5, $6, NOW(), NOW()), ($7, $8, $9, NOW(), NOW());
	`, 1, "Elementary School", 1, 2, "Middle School", 1, 3, "High School", 1)
	assert.NoError(t, err)

	therapistID := uuid.New()
	_, err = testDB.Exec(ctx, `
        INSERT INTO therapist (id, first_name, last_name, email, active, schools, district_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, therapistID, "District", "Therapist", "therapist@district.com", true, []int{1, 2, 3}, 1, time.Now(), time.Now())
	assert.NoError(t, err)

	testDOB := time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 2; j++ {
			_, err = testDB.Exec(ctx, `
                INSERT INTO student (id, first_name, last_name, dob, therapist_id, school_id, grade, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            `, uuid.New(), fmt.Sprintf("School%d", i), fmt.Sprintf("Student%d", j), testDOB, therapistID, i, i*2, time.Now(), time.Now())
			assert.NoError(t, err)
		}
	}

	students, err := repo.GetStudents(ctx, nil, PtrInt(1), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	students, err = repo.GetStudents(ctx, nil, PtrInt(2), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	students, err = repo.GetStudents(ctx, nil, PtrInt(3), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	students, err = repo.GetStudents(ctx, nil, nil, uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 6)

	students, err = repo.GetStudents(ctx, PtrInt(2), PtrInt(1), uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)
	for _, student := range students {
		assert.Equal(t, 2, *student.Grade)
		assert.Equal(t, 1, student.SchoolID)
	}
}
