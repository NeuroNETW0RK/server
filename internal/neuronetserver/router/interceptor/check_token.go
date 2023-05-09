package interceptor

import (
	"NeuroNET/internal/neuronetserver/model"
	"NeuroNET/internal/pkg/code"
	auth "NeuroNET/internal/pkg/jwt"
	"NeuroNET/internal/pkg/message"
	"NeuroNET/pkg/errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func (i *interceptor) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken := c.Request.Header.Get("Authorization")
		if rawToken == "" {
			message.Failed(c, errors.WithCode(code.ErrTokenInvalid, "Can not get Token in interceptor"))
			c.Abort()
			return
		}

		token := strings.Split(rawToken, " ")
		j := auth.NewJWT()
		claims, err := j.ParseToken(token[1])
		if err != nil {
			message.Failed(c, errors.WithCode(code.ErrTokenInvalid, "Can not parse Token in interceptor"))
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set(model.UserID, claims.ID)
		c.Next()
	}
}
