package app

import (
	"context"
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/ai/provider"
	"github.com/Veeru123456778/ai-interview-service/internal/interview"
	"github.com/Veeru123456778/ai-interview-service/internal/user"
	"github.com/Veeru123456778/ai-interview-service/internal/auth"
	"github.com/Veeru123456778/ai-interview-service/internal/health"
	"github.com/Veeru123456778/ai-interview-service/internal/middleware"
	"github.com/Veeru123456778/ai-interview-service/internal/resume"
	"github.com/gin-gonic/gin"
)

func NewRouter(application *Application) *gin.Engine {

	// Production mode
	if application.Config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(middleware.RequestID())

	router.Use(middleware.Logger(application.Logger))

	// Max resume upload size
	router.MaxMultipartMemory = int64(application.Config.Resume.MaxUploadSizeMB) << 20

	// Trust no proxy by default
	_ = router.SetTrustedProxies(nil)

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"service": application.Config.App.Name,
				"status":  "UP",
				"time":    time.Now().UTC(),
			},
		})
	})

	// -------------------------
	// Health Module
	// -------------------------
	healthService := health.NewService(
		application.DB,
		application.Redis,
		application.Config.App.Name,
		"1.0.0",
	)

	health.RegisterRoutes(router, healthService)

	// -------------------------
	// Resume Module
	// -------------------------

	resumeRepository := resume.NewRepository(application.DB)

	resumeExtractor := resume.NewExtractor()
	resumeNormalizer := resume.NewNormalizer()
	resumeBuilder := resume.NewIntelligenceBuilder()

	geminiClient, err := provider.NewClient(
		context.Background(),
		application.Config.Gemini.APIKey,
		application.Config.Gemini.Model,
		application.Config.Gemini.MaxRetries,
    )

	if err != nil {
		panic(err)
	}

	promptLoader := provider.NewPromptLoader()
	geminiProvider := provider.NewGeminiProvider(geminiClient, promptLoader)


	// -------------------------
	// User Module
	// -------------------------

	userRepository := user.NewRepository(application.DB)
	userService := user.NewService(userRepository)

	resumeParser := resume.NewParser(geminiProvider)

	resumeService := resume.NewService(
		resumeRepository,
		userService,
		resumeExtractor,
		resumeNormalizer,
		resumeParser,
		resumeBuilder,
	)


	// -------------------------
	// Interview Module
	// -------------------------

	interviewRepository := interview.NewRepository(application.DB)

	interviewService := interview.NewService(
		interviewRepository,
		resumeService,
		userService,
	)

	authService := auth.NewService(application.Config.Supabase.JWTSecret)

	api := router.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(authService))

	user.RegisterRoutes(api, userService)

	resume.RegisterRoutes(api, resumeService)

	interview.RegisterRoutes(api, interviewService)

	return router
}
