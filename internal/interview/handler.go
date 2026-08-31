package interview

import (
	"net/http"

	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ----------------------------------------------------------------------
// Register Routes
// ----------------------------------------------------------------------

func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {

	interviews := router.Group("/interviews")

	interviews.POST("", h.CreateInterview)
	interviews.GET("", h.ListInterviews)
	interviews.GET("/:interviewId", h.GetInterview)

	// Starts a new runtime interview session.
	interviews.POST("/:interviewId/start", h.StartInterview)
}

// ----------------------------------------------------------------------
// Create Interview
// ----------------------------------------------------------------------

func (h *Handler) CreateInterview(c *gin.Context) {

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			apperrors.ErrUnauthorized,
		)
		return
	}

	var request CreateInterviewRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"code":    "INVALID_REQUEST",
				"message": err.Error(),
			},
		)
		return
	}

	response, err := h.service.CreateInterview(
		c.Request.Context(),
		userID,
		&request,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ----------------------------------------------------------------------
// Get Interview
// ----------------------------------------------------------------------

func (h *Handler) GetInterview(c *gin.Context) {

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			apperrors.ErrUnauthorized,
		)
		return
	}

	response, err := h.service.GetInterview(
		c.Request.Context(),
		c.Param("interviewId"),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ----------------------------------------------------------------------
// List Interviews
// ----------------------------------------------------------------------

func (h *Handler) ListInterviews(c *gin.Context) {

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			apperrors.ErrUnauthorized,
		)
		return
	}

	response, err := h.service.ListInterviews(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ----------------------------------------------------------------------
// Start Interview
// ----------------------------------------------------------------------

func (h *Handler) StartInterview(c *gin.Context) {

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(
			http.StatusUnauthorized,
			apperrors.ErrUnauthorized,
		)
		return
	}

	state, err := h.service.StartInterview(
		c.Request.Context(),
		c.Param("interviewId"),
		userID,
	)
	if err != nil {
		writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, state)
}

// ----------------------------------------------------------------------
// Shared Error Writer
// ----------------------------------------------------------------------

func writeError(c *gin.Context, err error) {

	if appError, ok := err.(*apperrors.AppError); ok {
		c.JSON(appError.StatusCode, appError)
		return
	}

	c.JSON(
		http.StatusInternalServerError,
		apperrors.ErrInternalServer,
	)
}