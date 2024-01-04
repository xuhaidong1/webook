package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/xuhaidong1/go-generic-tools/pluginsx/logx"
	"github.com/xuhaidong1/go-generic-tools/pluginsx/saramax"
	"github.com/xuhaidong1/webook/internal/repository"
	"time"
)

const topicDownloadEvent = "app_download_event"

type DownloadEvent struct {
	Aid int64
	Uid int64
}

type Producer interface {
	ProduceDownloadEvent(evt DownloadEvent) error
}

var _ Consumer = &InteractiveDownloadConsumer{}

type InteractiveDownloadConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logx.Logger
}

func NewInteractiveReadEventConsumer(
	client sarama.Client,
	l logx.Logger,
	repo repository.InteractiveRepository) *InteractiveDownloadConsumer {
	ic := &InteractiveDownloadConsumer{
		repo:   repo,
		client: client,
		l:      l,
	}
	return ic
}

// Start 这边就是自己启动 goroutine 了
func (r *InteractiveDownloadConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{topicDownloadEvent},
			saramax.NewHandler[DownloadEvent](r.l, r.Consume))
		if er != nil {
			r.l.Error("退出了消费循环异常", logx.Error(er))
		}
	}()
	return err
}

func (r *InteractiveDownloadConsumer) StartBatch() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{topicDownloadEvent},
			saramax.NewBatchHandler[DownloadEvent](r.l, 100, r.BatchConsume))
		if er != nil {
			r.l.Error("退出了消费循环异常", logx.Error(er))
		}
	}()
	return err
}

func (r *InteractiveDownloadConsumer) Consume(msg *sarama.ConsumerMessage,
	evt DownloadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.repo.IncrReadCnt(ctx, "article", evt.Aid)
	return err
}

func (r *InteractiveDownloadConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	evts []DownloadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	bizs := make([]string, 0, len(msgs))
	ids := make([]int64, 0, len(msgs))
	for _, evt := range evts {
		bizs = append(bizs, "article")
		ids = append(ids, evt.Uid)
	}
	return r.repo.BatchIncrReadCnt(ctx, bizs, ids)
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{
		producer: producer,
	}
}

func (s *SaramaSyncProducer) ProduceDownloadEvent(evt DownloadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.
		SendMessage(&sarama.ProducerMessage{
			Topic: topicDownloadEvent,
			Key:   sarama.ByteEncoder(val),
		})
	return err
}
