package repository

import (
	"context"
	"github.com/xuhaidong1/webook/webook-be/internal/domain"
	"github.com/xuhaidong1/webook/webook-be/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	res, err := r.dao.FindByEmail(ctx, email)
	return domain.User{
		Email:    res.Email,
		Password: res.Password,
	}, err
}
