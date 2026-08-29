package health

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func RegisterRoutes(router *gin.Engine, service *Service) {
	handler := &Handler{
		service: service,
	}

	router.GET("/health", handler.Health)
	router.GET("/health/live", handler.Live)
	router.GET("/health/ready", handler.Ready)
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.service.Health(),
	})
}

func (h *Handler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    h.service.Live(),
	})
}

func (h *Handler) Ready(c *gin.Context) {
	response := h.service.Ready(context.Background())

	statusCode := http.StatusOK
	if response["status"] == "NOT_READY" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    response,
	})
}