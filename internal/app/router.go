package app

import (
	"time"

	"github.com/Veeru123456778/ai-interview-service/internal/health"
	"github.com/gin-gonic/gin"
)

func NewRouter(application *Application) *gin.Engine {

	// Production mode
	if application.Config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

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
	// Future Modules
	// -------------------------
	// auth.RegisterRoutes(router, ...)
	// resume.RegisterRoutes(router, ...)
	// interview.RegisterRoutes(router, ...)
	// user.RegisterRoutes(router, ...)
	// websocket.RegisterRoutes(router, ...)

	return router
}