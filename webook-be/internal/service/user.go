package service

import (
	"context"
	"errors"
	"github.com/xuhaidong1/webook/webook-be/internal/domain"
	"github.com/xuhaidong1/webook/webook-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
	ErrWrongPassword      = errors.New("用户名或密码错误")
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

func (svc UserService) Login(ctx context.Context, u domain.User) error {
	res, err := svc.repo.FindByEmail(ctx, u.Email)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(u.Password))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}
