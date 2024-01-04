package dao

import (
	"context"
	"gorm.io/gorm"
)

type AppDAO interface {
	GetById(ctx context.Context, id int64) (App, error)
}

type appDAO struct {
	db *gorm.DB
}

func (a *appDAO) GetById(ctx context.Context, id int64) (App, error) {
	//TODO implement me
	panic("implement me")
}

type App struct {
	Id   int64 `gorm:"primaryKey,autoIncrement"`
	Name string
	Url  string
}
