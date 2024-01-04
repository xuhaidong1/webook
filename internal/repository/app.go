package repository

import (
	"context"
	"github.com/xuhaidong1/webook/internal/domain"
	"github.com/xuhaidong1/webook/internal/repository/dao"
)

type AppRepository interface {
	GetById(ctx context.Context, id int64) (domain.App, error)
}

type appRepository struct {
	dao dao.AppDAO
}

func (a *appRepository) GetById(ctx context.Context, id int64) (domain.App, error) {
	res, err := a.dao.GetById(ctx, id)
	if err != nil {
		return domain.App{}, err
	}
	return domain.App{
		Id:   res.Id,
		Name: res.Name,
		Url:  res.Url,
	}, nil
}
