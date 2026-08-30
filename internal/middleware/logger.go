package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

 	"github.com/Veeru123456778/ai-interview-service/internal/auth"
)

func Logger(log *zap.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		duration := time.Since(start)

		requestID, _ := c.Get(RequestIDKey)

		fields := []zap.Field{
			zap.String("request_id", requestID.(string)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		}

		if user, ok := c.Get("authenticated_user"); ok {
			fields = append(fields, zap.String("user_id", user.(*auth.UserContext).SupabaseUserID))
		}

		log.Info("http_request", fields...)
	}
}