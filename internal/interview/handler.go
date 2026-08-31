package interview

import (
	"net/http"

	"github.com/Veeru123456778/ai-interview-service/internal/auth"
	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func RegisterRoutes(router *gin.RouterGroup, service Service) {
	handler := &Handler{
		service: service,
	}

	router.POST("/interviews", handler.CreateInterview)
	router.GET("/interviews", handler.ListInterviews)
	router.GET("/interviews/:interviewId", handler.GetInterview)
}

// ----------------------------------------------------------------------
// POST /api/v1/interviews
// ----------------------------------------------------------------------

func (h *Handler) CreateInterview(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	var request CreateInterviewRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, &apperrors.AppError{
			Code:       "INVALID_REQUEST",
			Message:    "Invalid interview request.",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	interviewResponse, err := h.service.CreateInterview(
		c.Request.Context(),
		user.SupabaseUserID,
		&request,
	)
	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr)
			return
		}

		response.InternalServer(c)
		return
	}

	response.Created(c, interviewResponse)
}

// ----------------------------------------------------------------------
// GET /api/v1/interviews/:interviewId
// ----------------------------------------------------------------------

func (h *Handler) GetInterview(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	interviewID := c.Param("interviewId")
	if interviewID == "" {
		response.Error(c, &apperrors.AppError{
			Code:       "INVALID_INTERVIEW_ID",
			Message:    "Interview ID is required.",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	interviewResponse, err := h.service.GetInterview(
		c.Request.Context(),
		interviewID,
		user.SupabaseUserID,
	)
	if err != nil {

		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr)
			return
		}

		response.InternalServer(c)
		return
	}

	response.OK(c, interviewResponse)
}

// ----------------------------------------------------------------------
// GET /api/v1/interviews
// ----------------------------------------------------------------------

func (h *Handler) ListInterviews(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	interviews, err := h.service.ListInterviews(
		c.Request.Context(),
		user.SupabaseUserID,
	)
	if err != nil {
		response.InternalServer(c)
		return
	}

	response.OK(c, interviews)
}