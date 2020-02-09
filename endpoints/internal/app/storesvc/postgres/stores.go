package postgres

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"

	"github.com/cage1016/mask/internal/app/model"
	"github.com/cage1016/mask/internal/pkg/errors"
)

var (
	ErrQueryStoreFromStoresDB   = errors.New("query store from DB failed")
	ErrInsertOrUpdateToStoresDB = errors.New("insert or update DB failed")
	ErrDeleteStoresDB           = errors.New("delete stores DB failed")
	ErrTxCommit                 = errors.New("Tx commit failed")
)

var _ model.StoreRepository = (*storeRepository)(nil)

type storeRepository struct {
	db  *sqlx.DB
	log log.Logger
}

// New instantiates a PostgreSQL implementation of givenEmail
// repository.
func New(db *sqlx.DB, log log.Logger) model.StoreRepository {
	return &storeRepository{db, log}
}

func (s storeRepository) Query(ctx context.Context, centerLng, centerLat, swLng, neLng, swLat, neLat float64, max uint64) ([]model.Store, error) {
	q := `SELECT *, point ($1, $2) <@> point(longitude, latitude)::point as distance
			FROM (select * from stores where longitude >= $3 and longitude <= $4 and latitude >= $5 and latitude <= $6) as a
			ORDER BY distance limit $7;`

	stores := []model.Store{}
	if err := s.db.SelectContext(ctx, &stores, q, centerLng, centerLat, swLng, neLng, swLat, neLat, max); err != nil {
		return stores, errors.Wrap(ErrQueryStoreFromStoresDB, err)
	}

	return stores, nil
}

func (s storeRepository) Insert(ctx context.Context, stores []model.Store) error {
	if _, err := s.db.Exec(`delete from stores`); err != nil {
		level.Error(s.log).Log("method", "s.db.Exec", "sql", "delete from stores", "err", err)
		return errors.Wrap(ErrDeleteStoresDB, err)
	}

	q := `INSERT INTO public.stores (id, name, phone, address, mask_adult, mask_child, available, note, longitude, latitude, custom_note, website, updated)
			VALUES (:id, :name, :phone, :address, :mask_adult, :mask_child, :available, :note, :longitude, :latitude, :custom_note, :website, :updated)
			ON CONFLICT (id)
				DO UPDATE
				SET name        = :name,
					phone       = :phone,
					address     = :address,
					mask_adult  = :mask_adult,
					mask_child  = :mask_child,
					available   = :available,
					note        = :note,
					longitude   = :longitude,
					latitude    = :latitude,
					custom_note = :custom_note,
					website     = :website,
					updated     = :updated`

	ctx = context.Background()
	tx := s.db.MustBeginTx(ctx, nil)

	for _, store := range stores {
		if _, err := tx.NamedExecContext(ctx, q, store); err != nil {
			level.Error(s.log).Log("method", "tx.NamedExecContext", "err", err)
			return errors.Wrap(ErrInsertOrUpdateToStoresDB, err)
		}
	}

	err := tx.Commit()
	if err != nil {
		return errors.Wrap(ErrTxCommit, err)
	}

	return nil
}
