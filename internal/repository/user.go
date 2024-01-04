package repository

import (
	"context"
	"github.com/xuhaidong1/webook/webook-be/internal/domain"
	"github.com/xuhaidong1/webook/webook-be/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, u domain.User) (domain.User, error) {
	res, err := r.dao.FindByEmail(ctx, u.Email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       res.Id,
		Email:    res.Email,
		Password: res.Password,
	}, nil
}

func (r *UserRepository) FindById(ctx context.Context, u domain.User) (domain.User, error) {
	res, err := r.dao.FindById(ctx, u.Id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       res.Id,
		Email:    res.Email,
		Password: res.Password,
		Profile: domain.Profile{
			Nickname:     res.Nickname,
			Birth:        res.Birth,
			Introduction: res.Introduction,
		},
	}, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, u domain.User) error {
	return r.dao.UpdateProfileById(ctx, u)
}
