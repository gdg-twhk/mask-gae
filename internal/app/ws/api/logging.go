package api

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/cage1016/mask/internal/app/ws"
)

var _ ws.Service = (*loggingMiddleware)(nil)

type loggingMiddleware struct {
	logger log.Logger
	svc    ws.Service
}

// LoggingMiddleware adds logging facilities to the adapter.
func LoggingMiddleware(svc ws.Service, logger log.Logger) ws.Service {
	return &loggingMiddleware{logger, svc}
}

func (lm *loggingMiddleware) Publish(ctx context.Context, topic string, msg string) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			level.Warn(lm.logger).Log("method", "Publish", "topic", topic, "err", err, "took", time.Since(begin))
			return
		}
	}(time.Now())

	return lm.svc.Publish(ctx, topic, msg)
}

func (lm *loggingMiddleware) Subscribe(topic string, channel *ws.Channel) (err error) {
	defer func(begin time.Time) {
		if err != nil {
			level.Warn(lm.logger).Log("method", "Subscribe", "topic", topic, "err", err, "took", time.Since(begin))
			return
		}
	}(time.Now())

	return lm.svc.Subscribe(topic, channel)
}
