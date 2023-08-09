package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/xuhaidong1/webook/webook-be/internal/domain"
	"github.com/xuhaidong1/webook/webook-be/internal/service"
	"net/http"
)

// UserHandler 定义与用户有关的路由
type UserHandler struct {
	emailExp *regexp.Regexp
	pwdExp   *regexp.Regexp
	svc      *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

		userIdKey = "userId"
	)
	//预编译正则表达式提高校验速度
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	pwdExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:      svc,
		emailExp: emailExp,
		pwdExp:   pwdExp,
	}
}

func (u *UserHandler) RegisterUserRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.Signup)
	server.POST("/users/edit", u.Edit)
	server.POST("/users/login", u.Login)
	server.POST("/users/profile", u.Profile)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	//作用域最小化
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	//Bind 根据content-type解析数据到req里面
	//解析错了会写回一个4xx的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//校验
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		//TODO 记录日志
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}

	ok, err = u.pwdExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}
	//todo
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "注册成功哇")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrWrongPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功哇")
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) Profile(ctx *gin.Context) {

}
