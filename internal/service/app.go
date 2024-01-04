package service

import (
	"context"
	"github.com/xuhaidong1/go-generic-tools/pluginsx/logx"
	"github.com/xuhaidong1/webook/internal/domain"
	"github.com/xuhaidong1/webook/internal/events"
	"github.com/xuhaidong1/webook/internal/repository"
)

type AppService interface {
	Download(ctx context.Context, id, uid int64) (domain.App, error)
}

type appService struct {
	repo     repository.AppRepository
	producer events.Producer
	logger   logx.Logger
}

func (svc *appService) Download(ctx context.Context, id, uid int64) (domain.App, error) {
	res, err := svc.repo.GetById(ctx, id)
	go func() {
		if err == nil {
			er := svc.producer.ProduceDownloadEvent(events.DownloadEvent{
				Aid: id,
				Uid: id,
			})
			if er != nil {
				svc.logger.Error("发送消息失败",
					logx.Int64("uid", uid),
					logx.Int64("aid", id),
					logx.Error(err))
			}
		}
	}()
	return res, err
}
