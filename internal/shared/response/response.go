package response

import (
	"net/http"

	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}

func Error(c *gin.Context, err *apperrors.AppError) {
	c.JSON(err.StatusCode, gin.H{
		"success": false,
		"error": gin.H{
			"code":    err.Code,
			"message": err.Message,
		},
	})
}

func InternalServer(c *gin.Context) {
	Error(c, apperrors.ErrInternalServer)
}

func OK(c *gin.Context, data any) {
	Success(c, http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	Success(c, http.StatusCreated, data)
}