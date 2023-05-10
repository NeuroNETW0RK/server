package interceptor

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"neuronet/pkg/log"
)

func (i *interceptor) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(log.KeyRequestID)

		if rid == "" {
			rid = uuid.Must(uuid.NewV4()).String()
			c.Request.Header.Set(log.KeyRequestID, rid)
			c.Set(log.KeyRequestID, rid)
		}

		// Set XRequestIDKey header
		c.Writer.Header().Set(log.KeyRequestID, rid)
		c.Next()
	}
}
