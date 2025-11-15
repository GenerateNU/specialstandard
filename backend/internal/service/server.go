package service

import (
	"log/slog"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"
	"specialstandard/internal/s3_client"
	"specialstandard/internal/service/handler/auth"
	"specialstandard/internal/service/handler/game_content"
	"specialstandard/internal/service/handler/game_result"
	"specialstandard/internal/service/handler/resource"
	"specialstandard/internal/service/handler/school"
	"specialstandard/internal/service/handler/session"
	"specialstandard/internal/service/handler/session_resource"
	sessionstudent "specialstandard/internal/service/handler/session_student"
	"specialstandard/internal/service/handler/student"
	"specialstandard/internal/service/handler/theme"
	"specialstandard/internal/service/handler/therapist"
	"specialstandard/internal/storage"
	"specialstandard/internal/storage/postgres"

	"context"
	"net/http"
	supabase_auth "specialstandard/internal/auth"

	"specialstandard/internal/service/handler/district"

	go_json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

type App struct {
	Server   *fiber.App
	Repo     *storage.Repository
	S3Bucket *s3_client.Client
}

// Initialize the App union type containing a fiber app, a repository, and a climatiq client.
func InitApp(config config.Config) *App {
	ctx := context.Background()
	repo := postgres.NewRepository(ctx, config.DB)
	bucket, err := s3_client.NewClient(config.S3Bucket)
	if err != nil {
		slog.Error("bucket cannot be configured")
	}

	app := SetupApp(config, repo, bucket)

	return &App{
		Server:   app,
		Repo:     repo,
		S3Bucket: bucket,
	}
}

// Setup the fiber app with the specified configuration, database, and S3 client.
func SetupApp(config config.Config, repo *storage.Repository, bucket *s3_client.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder:  go_json.Marshal,
		JSONDecoder:  go_json.Unmarshal,
		ErrorHandler: errs.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(favicon.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Use logging middleware
	app.Use(logger.New())

	// Use CORS middleware to configure CORS and handle preflight/OPTIONS requests.
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:3001,http://localhost:8080,http://127.0.0.1:8080,http://127.0.0.1:3000,https://clownfish-app-wq7as.ondigitalocean.app,https://king-prawn-app-n5vk6.ondigitalocean.app",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS", // Using these methods.
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // Allow cookies
		ExposeHeaders:    "Content-Length, X-Request-ID",
	}))

	app.Static("/api", "/app/api")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/api/openapi.yaml",
		DeepLinking: false,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Welcome to The Special Standard!")
	})

	apiV1 := app.Group("/api/v1")

	apiV1.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
	// Setup

	SupabaseAuthHandler := auth.NewHandler(config.Supabase, repo.Therapist)

	authGroup := apiV1.Group("/auth")
	authGroup.Post("/login", SupabaseAuthHandler.Login)
	authGroup.Post("/signup", SupabaseAuthHandler.SignUp)

	if !config.TestMode {
		apiV1.Use(supabase_auth.Middleware(&config.Supabase))
	} else {
		apiV1.Use(func(c *fiber.Ctx) error {
			c.Locals("user", "test-user")
			return c.Next()
		})
	}

	themeHandler := theme.NewHandler(repo.Theme)
	apiV1.Route("/themes", func(r fiber.Router) {
		r.Post("/", themeHandler.CreateTheme)
		r.Get("/", themeHandler.GetThemes)
		r.Get("/:id", themeHandler.GetThemeByID)
		r.Patch("/:id", themeHandler.PatchTheme)
		r.Delete("/:id", themeHandler.DeleteTheme)
	})

	therapistHandler := therapist.NewHandler(repo.Therapist)
	apiV1.Route("/therapists", func(r fiber.Router) {
		r.Get("/:id", therapistHandler.GetTherapistByID)
		r.Post("/", therapistHandler.CreateTherapist)
		r.Get("/", therapistHandler.GetTherapists)
		r.Delete("/:id", therapistHandler.DeleteTherapist)
		r.Patch("/:id", therapistHandler.PatchTherapist)
	})

	resourceHandler := resource.NewHandler(repo.Resource, bucket)
	apiV1.Route("/resources", func(r fiber.Router) {
		r.Post("/", resourceHandler.PostResource)
		r.Get("/", resourceHandler.GetResources)
		r.Get("/:id", resourceHandler.GetResourceByID)
		r.Patch("/:id", resourceHandler.UpdateResource)
		r.Delete("/:id", resourceHandler.DeleteResource)
	})

	sessionStudentHandler := sessionstudent.NewHandler(repo.SessionStudent)
	apiV1.Route("/session_students", func(r fiber.Router) {
		r.Post("/", sessionStudentHandler.CreateSessionStudent)
		r.Delete("/", sessionStudentHandler.DeleteSessionStudent)
		r.Patch("/", sessionStudentHandler.PatchStudentSessionRatings)
	})

	studentHandler := student.NewHandler(repo.Student)
	// Student route
	apiV1.Route("/students", func(r fiber.Router) {
		r.Get("/", studentHandler.GetStudents)
		r.Get("/:id", studentHandler.GetStudent)
		r.Delete("/:id", studentHandler.DeleteStudent)
		r.Post("/", studentHandler.AddStudent)
		r.Patch("/promote", studentHandler.PromoteStudents)
		r.Patch("/:id", studentHandler.UpdateStudent)
		r.Get("/:id/sessions", studentHandler.GetStudentSessions)
		r.Get("/:id/ratings", studentHandler.GetStudentRatings)
		r.Get("/:id/attendance", sessionStudentHandler.GetStudentAttendance)
	})

	sessionResourceHandler := session_resource.NewHandler(repo.SessionResource)
	apiV1.Route("/session-resource", func(r fiber.Router) {
		r.Post("/", sessionResourceHandler.PostSessionResource)
		r.Delete("/", sessionResourceHandler.DeleteSessionResource)
	})

	sessionHandler := session.NewHandler(repo.Session, repo.SessionStudent)

	apiV1.Route("/sessions", func(r fiber.Router) {
		r.Get("/", sessionHandler.GetSessions)
		r.Post("/", sessionHandler.PostSessions)
		r.Get("/:id", sessionHandler.GetSessionByID)
		r.Get("/:id/resources", sessionResourceHandler.GetSessionResources)
		r.Patch("/:id", sessionHandler.PatchSessions)
		r.Get("/:id/students", sessionHandler.GetSessionStudents)
		r.Delete("/:id", sessionHandler.DeleteSessions)
	})

	gameContentHandler := game_content.NewHandler(repo.GameContent)
	apiV1.Route("/game-contents", func(r fiber.Router) {
		r.Get("/", gameContentHandler.GetGameContents)
	})

	gameResultsHandler := game_result.NewHandler(repo.GameResult)
	apiV1.Route("/game-results", func(r fiber.Router) {
		r.Get("/", gameResultsHandler.GetGameResults)
		r.Post("/", gameResultsHandler.PostGameResult)
	})

	districtHandler := district.NewHandler(repo.District)
	apiV1.Route("/districts", func(r fiber.Router) {
		r.Get("/", districtHandler.GetDistricts)
		r.Get("/:id", districtHandler.GetDistrictByID)
	})

	schoolHandler := school.NewHandler(repo.School)
	apiV1.Route("/schools", func(r fiber.Router) {
		r.Get("/", schoolHandler.GetSchools)
	})
	

	// Handle 404 - Route not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
			"path":  c.Path(),
		})
	})

	app.Get("/secret", supabase_auth.Middleware(&config.Supabase), func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	return app
}
