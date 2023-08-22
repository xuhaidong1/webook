package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	//不需要登录校验的path
	Paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gob.Register(time.Now())
		for _, path := range l.Paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//if ctx.Request.URL.Path == "/users/login" || ctx.Request.URL.Path == "/users/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime := sess.Get("updateTime")
		now := time.Now().UnixMilli()
		if updateTime == nil || now-updateTime.(int64) > 10*1000 {
			updateTime = now
		}
		sess.Options(sessions.Options{
			//Secure: true,
			//HttpOnly:true,
			MaxAge: 10,
		})
		sess.Set("userId", id)
		sess.Set("updateTime", updateTime)
		err := sess.Save()
		if err != nil {
			return
		}
	}
}

// IgnorePath 链式调用 可配置path 不要对用户调用顺序有假设
func (l *LoginMiddlewareBuilder) IgnorePath(path string) *LoginMiddlewareBuilder {
	l.Paths = append(l.Paths, path)
	return l
}
