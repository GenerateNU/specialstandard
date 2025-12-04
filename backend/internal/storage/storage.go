package storage

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/dbinterface"
	"specialstandard/internal/storage/postgres/schema"
	"specialstandard/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsletterRepository interface {
	GetNewsletterByDate(ctx context.Context, date time.Time) (*models.Newsletter, error)
}

type SessionRepository interface {
	GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest, therapistid uuid.UUID) ([]models.Session, error)
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteRecurringSessions(ctx context.Context, id uuid.UUID) error
	PostSession(ctx context.Context, q dbinterface.Queryable, session *models.PostSessionInput) (*[]models.Session, error)
	PatchSession(ctx context.Context, id uuid.UUID, session *models.PatchSessionInput) (*models.Session, error)
	GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination, therapistId uuid.UUID) ([]models.SessionStudentsOutput, error)

	GetDB() *pgxpool.Pool
}

type SessionStudentRepository interface {
	CreateSessionStudent(ctx context.Context, q dbinterface.Queryable, input *models.CreateSessionStudentInput) (*[]models.SessionStudent, error)
	DeleteSessionStudent(ctx context.Context, input *models.DeleteSessionStudentInput) error
	RateStudentSession(ctx context.Context, input *models.PatchSessionStudentInput) (*models.SessionStudent, []models.SessionRating, error)
	GetStudentAttendance(ctx context.Context, params models.GetStudentAttendanceParams) (*int, *int, error)
	GetDB() *pgxpool.Pool
}

type StudentRepository interface {
	GetStudents(ctx context.Context, grade, schoolID *int, therapistID uuid.UUID, name string, pagination utils.Pagination) ([]models.Student, error)
	GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error)
	AddStudent(ctx context.Context, student models.Student) (models.Student, error)
	UpdateStudent(ctx context.Context, student models.Student) (models.Student, error)
	DeleteStudent(ctx context.Context, id uuid.UUID) error
	GetStudentSessions(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRepositoryRequest) ([]models.StudentSessionsOutput, error)
	GetStudentRatings(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRatingsRequest) ([]models.StudentSessionsWithRatingsOutput, error)
	PromoteStudents(ctx context.Context, input models.PromoteStudentsInput) error
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
	GetResources(ctx context.Context, themeID uuid.UUID, gradeLevel, resType, title, category, content, themeName string, week string, themeMonth, themeYear *int, pagination utils.Pagination) ([]models.ResourceWithTheme, error)
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

type GameContentRepository interface {
	GetGameContents(ctx context.Context, req models.GetGameContentRequest) ([]models.GameContent, error)
}

type GameResultRepository interface {
	GetGameResults(ctx context.Context, inputQuery *models.GetGameResultQuery, pagination utils.Pagination) ([]models.GameResult, error)
	PostGameResult(ctx context.Context, input models.PostGameResult) (*models.GameResult, error)
}

type DistrictRepository interface {
	GetDistricts(ctx context.Context) ([]models.District, error)
	GetDistrictByID(ctx context.Context, id int) (*models.District, error)
}

type SchoolRepository interface {
	GetSchools(ctx context.Context) ([]models.School, error)
	GetSchoolsByDistrict(ctx context.Context, districtID int) ([]models.School, error)
}

// VerificationRepository defines methods for verification code operations
type VerificationRepository interface {
	CreateVerificationCode(ctx context.Context, code models.VerificationCode) error
	VerifyCode(ctx context.Context, userID, code string) (bool, error)
}

type AuthRepository interface {
	GetUserEmail(ctx context.Context, userID string) (string, error)
	MarkEmailVerified(ctx context.Context, userID string) error
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
	GameContent     GameContentRepository
	GameResult      GameResultRepository
	District        DistrictRepository
	School          SchoolRepository
	Newsletter      NewsletterRepository
	Verification    VerificationRepository
	Auth            AuthRepository
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
		GameContent:     schema.NewGameContentRepository(db),
		GameResult:      schema.NewGameResultRepository(db),
		District:        schema.NewDistrictRepository(db),
		School:          schema.NewSchoolRepository(db),
		Newsletter:      schema.NewNewsletterRepository(db),
		Verification:    schema.NewVerificationRepository(db),
	}
}
