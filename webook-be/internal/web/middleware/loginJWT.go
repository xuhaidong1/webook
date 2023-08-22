package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuhaidong1/webook/webook-be/internal/web"
	"net/http"
	"strings"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	//不需要登录校验的path
	Paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//gob.Register(time.Now())
		for _, path := range l.Paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//与前端约定jwt放在header里的Authorization
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			//有人瞎搞
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		cliams := &web.UserClaims{}
		//ParseWithClaims 要传入cliams的指针，该函数会修改cliams
		token, err := jwt.ParseWithClaims(tokenStr, cliams, func(token *jwt.Token) (interface{}, error) {
			return []byte("eYaunJyLLsEO35a7zLDDzXG66O50tj5D"), nil
		})
		if err != nil {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//如果token过期了 会invalid
		if token == nil || !token.Valid || cliams.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if cliams.UserAgent != ctx.Request.UserAgent() {
			//安全问题 需要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		now := time.Now()
		//1min过期 10s刷新1次
		if cliams.ExpiresAt.Sub(now) < time.Second*50 {
			cliams.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			//设定过期时间需要重新生成token
			tokenStr, err := token.SignedString([]byte("eYaunJyLLsEO35a7zLDDzXG66O50tj5D"))
			if err != nil {
				//记录日志续约失败
				//return
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		//cliams放在ctx里 方便后面业务使用token，不需要再次解析token了
		ctx.Set("cliams", cliams)
	}
}

// IgnorePath 链式调用 可配置path 不要对用户调用顺序有假设
func (l *LoginJWTMiddlewareBuilder) IgnorePath(path string) *LoginJWTMiddlewareBuilder {
	l.Paths = append(l.Paths, path)
	return l
}
