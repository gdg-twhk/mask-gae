package postgres

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jmoiron/sqlx"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/errors"
)

var (
	ErrQueryStoreFromPharmaciesDB   = errors.New("query pharmacies from DB failed")
	ErrInsertOrUpdateToPharmaciesDB = errors.New("insert or update DB failed")
	ErrDeletePharmaciesDB           = errors.New("delete pharmacies DB failed")
	ErrTxCommit                     = errors.New("Tx commit failed")
)

var _ model.PharmacyRepository = (*pharmacyRepository)(nil)

type pharmacyRepository struct {
	db  *sqlx.DB
	log log.Logger
}

// New instantiates a PostgreSQL implementation of givenEmail
// repository.
func New(db *sqlx.DB, log log.Logger) model.PharmacyRepository {
	return &pharmacyRepository{db, log}
}

func (s pharmacyRepository) Query(ctx context.Context, centerLng, centerLat, swLng, neLng, swLat, neLat float64, max uint64) ([]model.Pharmacy, error) {
	q := `SELECT *, point ($1, $2) <@> point(longitude, latitude)::point as distance
			FROM (select * from pharmacies where longitude >= $3 and longitude <= $4 and latitude >= $5 and latitude <= $6) as a
			ORDER BY distance limit $7;`

	pharmacies := []model.Pharmacy{}
	if err := s.db.SelectContext(ctx, &pharmacies, q, centerLng, centerLat, swLng, neLng, swLat, neLat, max); err != nil {
		level.Error(s.log).Log("method", "s.db.SelectContext", "err", err)
		return pharmacies, errors.Wrap(ErrQueryStoreFromPharmaciesDB, err)
	}

	return pharmacies, nil
}

func (s pharmacyRepository) Insert(ctx context.Context, pharmacies []model.Pharmacy) error {
	if _, err := s.db.Exec(`delete from pharmacies`); err != nil {
		level.Error(s.log).Log("method", "s.db.Exec", "sql", "delete from pharmacies", "err", err)
		return errors.Wrap(ErrDeletePharmaciesDB, err)
	}

	q := `INSERT INTO pharmacies (id, name, phone, address, mask_adult, mask_child, available, note, longitude, latitude, custom_note, website, updated)
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

	for _, store := range pharmacies {
		if _, err := tx.NamedExecContext(ctx, q, store); err != nil {
			level.Error(s.log).Log("method", "tx.NamedExecContext", "err", err)
			return errors.Wrap(ErrInsertOrUpdateToPharmaciesDB, err)
		}
	}

	err := tx.Commit()
	if err != nil {
		return errors.Wrap(ErrTxCommit, err)
	}

	return nil
}
