package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/xuhaidong1/webook/webook-be/internal/domain"
	"github.com/xuhaidong1/webook/webook-be/internal/service"
	"net/http"
	"time"
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
	//	server.POST("/users/login", u.Login)
	server.POST("/users/login", u.LoginJWT)
	server.GET("/users/profile", u.ProfileJWT)
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
		ctx.String(http.StatusNotAcceptable, "请求解析错误")
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
		ctx.String(http.StatusNotAcceptable, "请求解析错误")
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrWrongUserorPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	//登录成功设置session
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly:true,
		MaxAge: 10,
	})
	sess.Set("userId", user.Id)
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功哇")
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusNotAcceptable, "请求解析错误")
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrWrongUserorPassword {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	//登录成功设置JWT
	cliams := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	tokenStr, err := token.SignedString([]byte("eYaunJyLLsEO35a7zLDDzXG66O50tj5D"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	ctx.String(http.StatusOK, "登录成功哇")
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName     string `json:"nickname"`
		Birth        string `json:"birth"`
		Introduction string `json:"introduction"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusNotAcceptable, "请求解析错误")
		return
	}
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	err := u.svc.Edit(ctx, domain.User{
		Id: userId.(int64),
		Profile: domain.Profile{
			Nickname:     req.NickName,
			Birth:        req.Birth,
			Introduction: req.Introduction,
		},
	})
	//todo 错误处理改造
	if err != nil {
		if err == service.ErrNickNameTooLong {
			ctx.String(http.StatusOK, "昵称过长")
		} else if err == service.ErrIntroductionTooLong {
			ctx.String(http.StatusOK, "简介过长")
		} else if err == service.ErrWrongBirthFormat {
			ctx.String(http.StatusOK, "生日格式错误")
		} else {
			ctx.String(http.StatusInternalServerError, "系统错误")
		}
		return
	}
	ctx.String(http.StatusOK, "修改成功")

}
func (u *UserHandler) Profile(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	profile, err := u.svc.Profile(ctx, domain.User{Id: userId.(int64)})
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}
	ctx.JSON(http.StatusOK, profile)
	ctx.String(http.StatusOK, "profileaaaaaa")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	//从jwt中取出userid
	c, _ := ctx.Get("cliams")
	cliams, ok := c.(*UserClaims) //1个返回值 类型不匹配会panic，2个返回值 类型不匹配会false
	if !ok {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	profile, err := u.svc.Profile(ctx, domain.User{Id: cliams.Uid})
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
	}
	ctx.JSON(http.StatusOK, profile)
	ctx.String(http.StatusOK, "profileaaaaaa")
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Uid int64
	// 自己随便加
	UserAgent string
}
