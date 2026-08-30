package middleware

import (
	"strings"

	"github.com/Veeru123456778/ai-interview-service/internal/auth"
	apperrors "github.com/Veeru123456778/ai-interview-service/internal/shared/errors"
	"github.com/Veeru123456778/ai-interview-service/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {

	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")

		user, err := authService.VerifyAccessToken(token)
		if err != nil {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		auth.SetUser(c, user)

		c.Next()
	}
}