package storage

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest) ([]models.Session, error)
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	PostSession(ctx context.Context, session *models.PostSessionInput) (*models.Session, error)
	PatchSession(ctx context.Context, id uuid.UUID, session *models.PatchSessionInput) (*models.Session, error)
	GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.SessionStudentsOutput, error)
}

type SessionStudentRepository interface {
	CreateSessionStudent(ctx context.Context, input *models.CreateSessionStudentInput) (*models.SessionStudent, error)
	DeleteSessionStudent(ctx context.Context, input *models.DeleteSessionStudentInput) error
	PatchSessionStudent(ctx context.Context, input *models.PatchSessionStudentInput) (*models.SessionStudent, error)
}

type StudentRepository interface {
	GetStudents(ctx context.Context, grade string, therapistID uuid.UUID, name string, pagination utils.Pagination) ([]models.Student, error)
	GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error)
	AddStudent(ctx context.Context, student models.Student) (models.Student, error)
	UpdateStudent(ctx context.Context, student models.Student) (models.Student, error)
	DeleteStudent(ctx context.Context, id uuid.UUID) error
	GetStudentSessions(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRepositoryRequest) ([]models.StudentSessionsOutput, error)
}

type ThemeRepository interface {
	CreateTheme(ctx context.Context, theme *models.CreateThemeInput) (*models.Theme, error)
	GetThemes(ctx context.Context, pagination utils.Pagination, filter *models.ThemeFilter) ([]models.Theme, error)
	GetThemeByID(ctx context.Context, id uuid.UUID) (*models.Theme, error)
	PatchTheme(ctx context.Context, id uuid.UUID, theme *models.UpdateThemeInput) (*models.Theme, error)
	DeleteTheme(ctx context.Context, id uuid.UUID) error
}

type TherapistRepository interface {
	GetTherapistByID(ctx context.Context, therapistID string) (*models.Therapist, error)
	GetTherapists(ctx context.Context, pagination utils.Pagination) ([]models.Therapist, error)
	CreateTherapist(ctx context.Context, therapist *models.CreateTherapistInput) (*models.Therapist, error)
	DeleteTherapist(ctx context.Context, therapistID string) error
	PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error)
}

type ResourceRepository interface {
	GetResources(ctx context.Context, themeID uuid.UUID, gradeLevel, resType, title, category, content, themeName string, date *time.Time, themeMonth, themeYear *int, pagination utils.Pagination) ([]models.ResourceWithTheme, error)
	GetResourceByID(ctx context.Context, id uuid.UUID) (*models.ResourceWithTheme, error)
	UpdateResource(ctx context.Context, id uuid.UUID, resourceBody models.UpdateResourceBody) (*models.Resource, error)
	CreateResource(ctx context.Context, resourceBody models.ResourceBody) (*models.Resource, error)
	DeleteResource(ctx context.Context, id uuid.UUID) error
}

type SessionResourceRepository interface {
	PostSessionResource(ctx context.Context, sessionResource models.CreateSessionResource) (*models.SessionResource, error)
	DeleteSessionResource(ctx context.Context, sessionResource models.DeleteSessionResource) error
	GetResourcesBySessionID(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.Resource, error)
}

type Repository struct {
	Resource        ResourceRepository
	db              *pgxpool.Pool
	Session         SessionRepository
	Student         StudentRepository
	Theme           ThemeRepository
	Therapist       TherapistRepository
	SessionStudent  SessionStudentRepository
	SessionResource SessionResourceRepository
}

func (r *Repository) Close() error {
	r.db.Close()
	return nil
}

func (r *Repository) GetDB() *pgxpool.Pool {
	return r.db
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:              db,
		Resource:        schema.NewResourceRepository(db),
		Session:         schema.NewSessionRepository(db),
		Student:         schema.NewStudentRepository(db),
		Theme:           schema.NewThemeRepository(db),
		Therapist:       schema.NewTherapistRepository(db),
		SessionStudent:  schema.NewSessionStudentRepository(db),
		SessionResource: schema.NewSessionResourceRepository(db),
	}
}
