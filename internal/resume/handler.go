package resume

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

	router.POST("/resumes", handler.UploadResume)
	router.GET("/resumes", handler.ListResumes)
	router.GET("/resumes/:resumeId", handler.GetResume)
}

// ----------------------------------------------------------------------
// POST /api/v1/resumes
// ----------------------------------------------------------------------

func (h *Handler) UploadResume(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	uploadedFile, err := c.FormFile("resume")
	if err != nil {
		response.Error(c, apperrors.ErrResumeRequired)
		return
	}

	file, err := uploadedFile.Open()
	if err != nil {
		response.InternalServer(c)
		return
	}
	defer file.Close()

	uploadResponse, err := h.service.CreateResume(
		c.Request.Context(),
		user.SupabaseUserID,
		file,
		uploadedFile.Filename,
	)

	if err != nil {
		if appErr, ok := err.(*apperrors.AppError); ok {
			response.Error(c, appErr)
			return
		}

		response.InternalServer(c)
		return
	}

	response.Created(c, uploadResponse)
}

// ----------------------------------------------------------------------
// GET /api/v1/resumes/:resumeId
// ----------------------------------------------------------------------

func (h *Handler) GetResume(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	resumeID := c.Param("resumeId")
	if resumeID == "" {
		response.Error(c, &apperrors.AppError{
			Code:       "INVALID_RESUME_ID",
			Message:    "Resume ID is required.",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	resumeResponse, err := h.service.GetResume(
		c.Request.Context(),
		resumeID,
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

	response.OK(c, resumeResponse)
}

// ----------------------------------------------------------------------
// GET /api/v1/resumes
// ----------------------------------------------------------------------

func (h *Handler) ListResumes(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	resumes, err := h.service.ListResumes(
		c.Request.Context(),
		user.SupabaseUserID,
	)

	if err != nil {
		response.InternalServer(c)
		return
	}

	response.OK(c, resumes)
}