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
    `, studentID1, "John", "Doe", testDOB, therapistID1, "5th Grade", "IEP Goals: Speech articulation", time.Now(), time.Now())
	assert.NoError(t, err)

	studentID2 := uuid.New()
	_, err = testDB.Pool.Exec(ctx, `
        INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `, studentID2, "Jane", "Smith", testDOB, therapistID2, "3rd Grade", "IEP Goals: Reading", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test 1: Get all students (no filters)
	students, err := repo.GetStudents(ctx, "", uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 2)

	// Test 2: Filter by grade
	students, err = repo.GetStudents(ctx, "5th Grade", uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Doe", students[0].LastName)
	assert.Equal(t, studentID1, students[0].ID)

	// Test 3: Filter by therapist
	students, err = repo.GetStudents(ctx, "", therapistID2, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)
	assert.Equal(t, therapistID2, students[0].TherapistID)

	// Test 4: Filter by name (first name)
	students, err = repo.GetStudents(ctx, "", uuid.Nil, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)

	// Test 5: Filter by name (last name)
	students, err = repo.GetStudents(ctx, "", uuid.Nil, "Smith", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "Smith", students[0].LastName)

	// Test 6: Multiple filters
	students, err = repo.GetStudents(ctx, "5th Grade", therapistID1, "John", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 1)
	assert.Equal(t, "John", students[0].FirstName)
	assert.Equal(t, "5th Grade", *students[0].Grade)

	// Test 7: Filter that returns no results
	students, err = repo.GetStudents(ctx, "NonexistentGrade", uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 0)

	// More Tests for Pagination Behaviour
	for i := 3; i <= 6; i++ {
		testDOB, _ := time.Parse("2006-01-02", "2004-09-24")
		_, err := testDB.Pool.Exec(ctx, `
            INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            `, uuid.New(), "Student", fmt.Sprintf("Number%d", i), testDOB, therapistID1, fmt.Sprintf("Grade %d", i), "IEP: GOALS!!", time.Now(), time.Now())
		assert.NoError(t, err)
	}

	// Test 8: Pagination - get all students
	students, err = repo.GetStudents(ctx, "", uuid.Nil, "", utils.NewPagination())
	assert.NoError(t, err)
	assert.Len(t, students, 6) // 2 original + 4 new

	// Test 9: Pagination - second page
	students, err = repo.GetStudents(ctx, "", uuid.Nil, "", utils.Pagination{
		Page:  2,
		Limit: 5,
	})
	assert.NoError(t, err)
	assert.Len(t, students, 1)
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
    `, studentID, "Jane", "Smith", testDOB, therapistID, "3rd Grade", "IEP Goals: Language comprehension", time.Now(), time.Now())
	assert.NoError(t, err)

	// Test
	student, err := repo.GetStudent(ctx, studentID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Smith", student.LastName)
	assert.Equal(t, "Jane", student.FirstName)
	assert.Equal(t, studentID, student.ID)
	assert.Equal(t, therapistID, student.TherapistID)
	assert.Equal(t, "3rd Grade", *student.Grade)
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
		Grade:       ptrString("2nd Grade"),
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
    `, studentID, "Sam", "Wilson", testDOB, therapistID, "4th Grade", "Initial IEP", time.Now(), time.Now())
	assert.NoError(t, err)

	// Update student data
	updatedStudent := models.Student{
		ID:          studentID,
		FirstName:   "Sam",
		LastName:    "Wilson-Updated",
		DOB:         ptrTime(testDOB),
		TherapistID: therapistID,
		Grade:       ptrString("5th Grade"),
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
	assert.Equal(t, "5th Grade", *verifyStudent.Grade)
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
    `, studentID, "Chris", "Brown", testDOB, therapistID, "6th Grade", "IEP Goals: Social communication", time.Now(), time.Now())
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
