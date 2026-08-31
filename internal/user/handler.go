package user

import (
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

	router.POST("/users/sync", handler.SyncUser)
	router.GET("/users/me", handler.GetMe)
}

// ----------------------------------------------------------------------
// POST /api/v1/users/sync
// ----------------------------------------------------------------------

func (h *Handler) SyncUser(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	syncResponse, err := h.service.SyncUser(
		c.Request.Context(),
		user.SupabaseUserID,
		user.Email,
	)
	if err != nil {
		response.InternalServer(c)
		return
	}

	response.OK(c, syncResponse)
}

// ----------------------------------------------------------------------
// GET /api/v1/users/me
// ----------------------------------------------------------------------

func (h *Handler) GetMe(c *gin.Context) {

	user, ok := auth.GetUser(c)
	if !ok {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	syncResponse, err := h.service.SyncUser(
		c.Request.Context(),
		user.SupabaseUserID,
		user.Email,
	)
	if err != nil {
		response.InternalServer(c)
		return
	}

	response.OK(c, syncResponse.User)
}