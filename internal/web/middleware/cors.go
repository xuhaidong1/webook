package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type CorsMiddlewareBuilder struct {
}

func NewCorsMiddlewareBuilder() *CorsMiddlewareBuilder {
	return &CorsMiddlewareBuilder{}
}

func (c *CorsMiddlewareBuilder) Build() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//是否允许cookie之类的东西
		AllowCredentials: true,
		//允许给前端的东西
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "xxx.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
