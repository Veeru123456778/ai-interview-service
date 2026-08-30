package middleware

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {

	return func(c *gin.Context) {

		requestID := uuid.NewString()

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}