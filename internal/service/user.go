package service

import (
	"context"
	"errors"
	"github.com/xuhaidong1/webook/webook-be/internal/domain"
	"github.com/xuhaidong1/webook/webook-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUserDuplicateEmail  = repository.ErrUserDuplicateEmail
	ErrWrongUserorPassword = errors.New("用户名或密码错误")
	ErrNickNameTooLong     = errors.New("昵称过长")
	ErrWrongBirthFormat    = errors.New("生日填写错误")
	ErrIntroductionTooLong = errors.New("简介过长")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(encrypted)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	res, err := svc.repo.FindByEmail(ctx, u)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrWrongUserorPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(u.Password))
	if err != nil {
		return domain.User{}, ErrWrongUserorPassword
	}
	return res, nil
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	if err := svc.CheckProfileInput(u); err != nil {
		return err
	}
	if err := svc.repo.UpdateProfile(ctx, u); err != nil {
		return err
	}
	return nil
}

func (svc *UserService) Profile(ctx context.Context, u domain.User) (domain.Profile, error) {
	res, err := svc.repo.FindById(ctx, u)
	if err != nil {
		return domain.Profile{}, err
	}
	return res.Profile, nil
}

func (svc *UserService) CheckProfileInput(u domain.User) error {
	if len(u.Nickname) > 50 {
		return ErrNickNameTooLong
	}
	if len(u.Introduction) > 1000 {
		return ErrIntroductionTooLong
	}
	//没设定生日没事
	if len(u.Birth) == 0 {
		return nil
	}
	const layout = "2006-01-01"
	_, err := time.Parse(layout, u.Birth)
	if err != nil {
		return ErrWrongBirthFormat
	}
	return nil
}
