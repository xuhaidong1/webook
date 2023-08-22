package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/xuhaidong1/webook/webook-be/config"
	"github.com/xuhaidong1/webook/webook-be/internal/repository"
	"github.com/xuhaidong1/webook/webook-be/internal/repository/dao"
	"github.com/xuhaidong1/webook/webook-be/internal/service"
	"github.com/xuhaidong1/webook/webook-be/internal/web"
	"github.com/xuhaidong1/webook/webook-be/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GOOS=linux GOARCH=arm go build -o weboook .
func main() {
	db := initDB()
	u := initUser(db)
	initRedis()
	server := initServer()
	u.RegisterUserRoutes(server)
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello, world")
	//})
	server.Run(":" + config.Config.Port)
}

func initServer() *gin.Engine {
	server := gin.Default()
	//添加跨域中间件
	server.Use(middleware.NewCorsMiddlewareBuilder().Build())
	//初始化session中间件
	//store := cookie.NewStore([]byte("secret"))
	////store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	////	[]byte("HWwqUEQuKi4nEfRuA6TjwabYT6iOZ8y3"), []byte("eYaunJyLLsEO35a7zLDDzXG66O50tj5D"))
	//////size最大空闲 连接数，authentication 身份验证 encryption 数据加密
	////if err != nil {
	////	panic(err)
	////}
	//server.Use(sessions.Sessions("mySession", store))
	//初始化登录校验中间件
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePath("/users/login").IgnorePath("/users/signup").Build())
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePath("/users/login").IgnorePath("/users/signup").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		//只在初始化过程中panic
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	cfg := config.Config.Redis
	cmd := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	return cmd
}
