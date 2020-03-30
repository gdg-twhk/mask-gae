package service

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/level"
)

type loggingMiddleware struct {
	logger log.Logger      `json:""`
	next   PharmacyService `json:""`
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next PharmacyService) PharmacyService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Pharmacy, err error) {
	defer func() {
		lm.logger.Log("method", "Query", "centerLng", centerLng, "centerLat", centerLat, "neLng", neLng, "neLat", neLat, "seLng", seLng, "seLat", seLat, "swLng", swLng, "swLat", swLat, "nwLng", nwLng, "nwLat", nwLat, "max", max, "err", err)
	}()

	return lm.next.Query(ctx, centerLng, centerLat, neLng, neLat, seLng, seLat, swLng, swLat, nwLng, nwLat, max)
}

func (lm loggingMiddleware) Sync(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "Sync", "err", err)
	}()

	return lm.next.Sync(ctx)
}

func (lm loggingMiddleware) SyncHandler(ctx context.Context, queueName string, taskName string) (err error) {
	defer func() {
		lm.logger.Log("method", "SyncHandler", "queueName", queueName, "taskName", taskName, "err", err)
	}()

	return lm.next.SyncHandler(ctx, queueName, taskName)
}

func (lm loggingMiddleware) FootGun(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "FootGun", "err", err)
	}()

	return lm.next.FootGun(ctx)
}

func (lm loggingMiddleware) HealthCheck(ctx context.Context) (updated string, err error) {
	defer func() {
		lm.logger.Log("method", "HealthCheck", "err", err)
	}()

	return lm.next.HealthCheck(ctx)
}
