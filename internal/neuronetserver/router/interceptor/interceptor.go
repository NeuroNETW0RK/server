package interceptor

import (
	"github.com/gin-gonic/gin"
)

type Interceptor interface {
	RequestID() gin.HandlerFunc
	JWTAuth() gin.HandlerFunc
}

type interceptor struct {
}

func New() Interceptor {
	return &interceptor{}
}
