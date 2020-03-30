package service

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/level"
)

type loggingMiddleware struct {
	logger log.Logger
	next   PharmacyService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next PharmacyService) PharmacyService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) TickerUpdate(ctx context.Context) (err error) {
	defer func() {
		lm.logger.Log("method", "TickerUpdate", "err", err)
	}()

	return lm.next.TickerUpdate(ctx)
}

func (lm loggingMiddleware) Query(ctx context.Context, centerLng float64, centerLat float64, neLng float64, neLat float64, seLng float64, seLat float64, swLng float64, swLat float64, nwLng float64, nwLat float64, max uint64) (items []model.Pharmacy, err error) {
	defer func() {
		lm.logger.Log("method", "Query", "centerLng", centerLng, "centerLat", centerLat, "neLng", neLng, "neLat", neLat, "seLng", seLng, "seLat", seLat, "swLng", swLng, "swLat", swLat, "nwLng", nwLng, "nwLat", nwLat, "max", max, "err", err)
	}()

	return lm.next.Query(ctx, centerLng, centerLat, neLng, neLat, seLng, seLat, swLng, swLat, nwLng, nwLat, max)
}
