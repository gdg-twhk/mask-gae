package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gomurphyx/sqlx"

	"github.com/cage1016/mask/internal/pkg/util"
	"github.com/cage1016/mask/internal/app/feedback/model"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/level"
)

var _ model.FeedbackRepository = (*feedbackRepository)(nil)

var (
	ErrInsertOrUpdateToFeedbackDB = errors.New("insert or update DB failed")
)

type feedbackRepository struct {
	db  *sqlx.DB
	log log.Logger
}

// New instantiates a PostgreSQL implementation of givenEmail
// repository.
func New(db *sqlx.DB, log log.Logger) model.FeedbackRepository {
	return &feedbackRepository{db, log}
}

func (f feedbackRepository) Insert(ctx context.Context, feedback model.Feedback) (string, error) {
	nt := fmt.Sprintf("feedback_%s", time.Now().In(util.Location).Format("2006_0102"))
	level.Info(f.log).Log("table", nt)

	err := f.GetLatestFeedbackTableName(ctx, nt)
	if err != nil {
		return "", err
	}

	q := fmt.Sprintf(`INSERT INTO public.%s (id, user_id, pharmacy_id, option_id, description, longitude, latitude)
						VALUES (:id, :user_id, :pharmacy_id, :option_id, :description, :longitude, :latitude);`, nt)
	if _, err := f.db.NamedExecContext(ctx, q, feedback); err != nil {
		level.Error(f.log).Log("method", "s.db.NamedExecContext", "err", err)
		return "", errors.Wrap(ErrInsertOrUpdateToFeedbackDB, err)
	}
	return feedback.ID, nil
}

func (f feedbackRepository) RetrieveByUserID(ctx context.Context, userID string, date string, offset uint64, limit uint64) (model.FeedbackItemPage, error) {
	lt := struct {
		Exists bool `db:"exists"`
	}{}

	nt := fmt.Sprintf("feedback_%s", date)
	q := `SELECT EXISTS(SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = $1);`
	if err := f.db.GetContext(ctx, &lt, q, nt); err != nil {
		level.Error(f.log).Log("method", "f.db.GetContext", "sql", q, "nt", nt, "err", err)
		return model.FeedbackItemPage{}, err
	}

	if !lt.Exists {
		return model.FeedbackItemPage{Items: []model.Feedback{}}, nil
	}

	items := []model.Feedback{}
	q = fmt.Sprintf(`select * from %s where user_id = $1 order by created_at desc limit $2 offset $3`, nt)
	if err := f.db.SelectContext(ctx, &items, q, userID, limit, offset); err != nil {
		level.Error(f.log).Log("method", "f.db.SelectContext", "sql", q, "userID", userID, "limit", limit, "offset", offset)
		return model.FeedbackItemPage{Items: []model.Feedback{}}, err
	}

	cq := fmt.Sprintf(`select count(*) from %s where user_id = :user_id	`, nt)
	total, err := total(ctx, f.db, cq, map[string]interface{}{
		"user_id": userID,
	})
	if err != nil {
		return model.FeedbackItemPage{Items: []model.Feedback{}}, err
	}

	return model.FeedbackItemPage{
		Items: items,
		PageMetadata: model.PageMetadata{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func (f feedbackRepository) RetrieveByPharmacyID(ctx context.Context, pharmacyID string, date string, offset uint64, limit uint64) (model.FeedbackItemPage, error) {
	lt := struct {
		Exists bool `db:"exists"`
	}{}

	nt := fmt.Sprintf("feedback_%s", date)
	q := `SELECT EXISTS(SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = $1);`
	if err := f.db.GetContext(ctx, &lt, q, nt); err != nil {
		level.Error(f.log).Log("method", "f.db.GetContext", "sql", q, "nt", nt, "err", err)
		return model.FeedbackItemPage{}, err
	}

	if !lt.Exists {
		return model.FeedbackItemPage{Items: []model.Feedback{}}, nil
	}

	items := []model.Feedback{}
	q = fmt.Sprintf(`select * from %s where pharmacy_id = $1 order by created_at desc limit $2 offset $3`, nt)
	if err := f.db.SelectContext(ctx, &items, q, pharmacyID, limit, offset); err != nil {
		level.Error(f.log).Log("method", "f.db.SelectContext", "sql", q, "pharmacyID", pharmacyID, "limit", limit, "offset", offset)
		return model.FeedbackItemPage{Items: []model.Feedback{}}, err
	}

	cq := fmt.Sprintf(`select count(*) from %s where pharmacy_id = :pharmacy_id`, nt)
	total, err := total(ctx, f.db, cq, map[string]interface{}{
		"pharmacy_id": pharmacyID,
	})
	if err != nil {
		return model.FeedbackItemPage{Items: []model.Feedback{}}, err
	}

	return model.FeedbackItemPage{
		Items: items,
		PageMetadata: model.PageMetadata{
			Total:  total,
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func (f feedbackRepository) ListOption(ctx context.Context) ([]model.Option, error) {
	var options []model.Option

	if err := f.db.SelectContext(ctx, &options, `select * from options`); err != nil {
		return options, err
	}
	return options, nil
}

func (f feedbackRepository) GetLatestFeedbackTableName(ctx context.Context, nt string) error {
	lt := struct {
		Exists bool `db:"exists"`
	}{}

	q := `SELECT EXISTS(SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = $1);`
	if err := f.db.GetContext(ctx, &lt, q, nt); err != nil {
		level.Error(f.log).Log("method", "f.db.GetContext", "sql", q, "nt", nt, "err", err)
		return err
	}

	if lt.Exists {
		return nil
	}

	q = fmt.Sprintf(`create table if not exists %s
		(
			id   varchar(21)  not null
				constraint %s_pkey
					primary key,
			user_id varchar(30) default ''::character varying not null,
			pharmacy_id varchar(10) default ''::character varying not null,
			option_id varchar(21) default ''::character varying not null,
			description varchar(1024) default ''::character varying not null,
			longitude double precision default 0.0 not null,
			latitude double precision default 0.0 not null,
			created_at timestamp with time zone default now() not null
		);

		alter table %s owner to postgres;`, nt, nt, nt)

	if _, err := f.db.ExecContext(ctx, q); err != nil {
		level.Error(f.log).Log("method", "f.db.ExecContext", "sql", q, "err", err)
		return err
	}

	return nil
}

func total(ctx context.Context, db *sqlx.DB, query string, params map[string]interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}
