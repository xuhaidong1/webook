package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/xuhaidong1/webook/webook-be/internal/repository"
	"github.com/xuhaidong1/webook/webook-be/internal/repository/dao"
	"github.com/xuhaidong1/webook/webook-be/internal/service"
	"github.com/xuhaidong1/webook/webook-be/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	u := initUser(db)
	server := initServer()
	u.RegisterUserRoutes(server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initServer() *gin.Engine {
	server := gin.Default()
	//添加跨域中间件
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//是否允许cookie之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "xxx.com")
		},
		MaxAge: 12 * time.Hour,
	}))
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
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook-be"))
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
