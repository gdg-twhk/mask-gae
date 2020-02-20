package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/cage1016/mask/internal/app/feedback/model"
)

type loggingMiddleware struct {
	logger log.Logger
	next   FeedbacksvcService
}

// LoggingMiddleware takes a logger as a dependency
// and returns a ServiceMiddleware.
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next FeedbacksvcService) FeedbacksvcService {
		return loggingMiddleware{level.Info(logger), next}
	}
}

func (lm loggingMiddleware) Options(ctx context.Context) (items []model.Option, err error) {
	defer func() {
		lm.logger.Log("method", "Options", "err", err)
	}()

	return lm.next.Options(ctx)
}

func (lm loggingMiddleware) PharmacyFeedBacks(ctx context.Context, pharmacyID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error) {
	defer func() {
		lm.logger.Log("method", "PharmacyFeedBacks", "pharmacyID", pharmacyID, "date", date, "offset", offset, "limit", limit, "err", err)
	}()

	return lm.next.PharmacyFeedBacks(ctx, pharmacyID, date, offset, limit)
}

func (lm loggingMiddleware) UserFeedBacks(ctx context.Context, userID string, date string, offset, limit uint64) (res model.FeedbackItemPage, err error) {
	defer func() {
		lm.logger.Log("method", "UserFeedBacks", "userID", userID, "date", date, "offset", offset, "limit", limit, "err", err)
	}()

	return lm.next.UserFeedBacks(ctx, userID, date, offset, limit)
}

func (lm loggingMiddleware) InsertFeedBack(ctx context.Context, feedback model.Feedback) (id string, err error) {
	defer func() {
		lm.logger.Log("method", "InsertFeedBack", "feedback", feedback, "err", err)
	}()

	return lm.next.InsertFeedBack(ctx, feedback)
}
