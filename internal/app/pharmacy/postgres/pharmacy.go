package postgres

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/level"
)

var (
	ErrQueryStoreFromPharmaciesDB      = errors.New("query pharmacies from DB failed")
	ErrInsertOrUpdateToPharmaciesDB    = errors.New("insert or update DB failed")
	ErrExecContextPharmaciesDB         = errors.New("exec context pharmacies DB failed")
	ErrSQLAdvisoryLockFromPharmaciesDB = errors.New("psql advisory lock from from DB failed")
	ErrTxCommit                        = errors.New("Tx commit failed")
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
	ctx = context.Background()
	tx := s.db.MustBeginTx(ctx, nil)

	pharmacies := []model.Pharmacy{}
	if _, err := tx.ExecContext(ctx, `select pg_advisory_xact_lock(0);`); err != nil {
		tx.Rollback()
		level.Error(s.log).Log("method", "tx.ExecContext", "sql", "select pg_advisory_xact_lock(0);", "err", err)
		return pharmacies, errors.Wrap(ErrSQLAdvisoryLockFromPharmaciesDB, err)
	}

	q := `SELECT *, point ($1, $2) <@> point(longitude, latitude)::point as distance
			FROM (select * from pharmacies where longitude >= $3 and longitude <= $4 and latitude >= $5 and latitude <= $6) as a
			ORDER BY distance limit $7;`

	if err := tx.SelectContext(ctx, &pharmacies, q, centerLng, centerLat, swLng, neLng, swLat, neLat, max); err != nil {
		tx.Rollback()
		level.Error(s.log).Log("method", "s.db.SelectContext", "err", err)
		return pharmacies, errors.Wrap(ErrQueryStoreFromPharmaciesDB, err)
	}

	err := tx.Commit()
	if err != nil {
		level.Error(s.log).Log("method", "tx.Commit", "err", err)
		return pharmacies, errors.Wrap(ErrTxCommit, err)
	}

	return pharmacies, nil
}

func (s pharmacyRepository) Insert(ctx context.Context, pharmacies []model.Pharmacy) error {
	ctx = context.Background()
	tx := s.db.MustBeginTx(ctx, nil)

	if _, err := tx.ExecContext(ctx, `create table t as select * from pharmacies with no data;`); err != nil {
		tx.Rollback()
		level.Error(s.log).Log("method", "tx.ExecContext", "sql", "create table T as select * from pharmacies with no data;", "err", err)
		return errors.Wrap(ErrExecContextPharmaciesDB, err)
	}

	q := `INSERT INTO T (id, name, phone, address, mask_adult, mask_child, available, note, longitude, latitude,
						   custom_note, website, updated, service_periods, service_note, county, town, cunli)
			VALUES (:id, :name, :phone, :address, :mask_adult, :mask_child, :available, :note, :longitude, :latitude, :custom_note,
					:website, :updated, :service_periods, :service_note, :county, :town, :cunli);`

	for _, store := range pharmacies {
		if _, err := tx.NamedExecContext(ctx, q, store); err != nil {
			tx.Rollback()
			level.Error(s.log).Log("method", "tx.NamedExecContext", "err", err)
			return errors.Wrap(ErrInsertOrUpdateToPharmaciesDB, err)
		}
	}

	if _, err := tx.ExecContext(ctx, `select pg_advisory_xact_lock(0);`); err != nil {
		tx.Rollback()
		level.Error(s.log).Log("method", "tx.ExecContext", "sql", "select pg_advisory_xact_lock(0);", "err", err)
		return errors.Wrap(ErrSQLAdvisoryLockFromPharmaciesDB, err)
	}

	if _, err := tx.ExecContext(ctx, `select pg_advisory_xact_lock(0); drop table pharmacies; alter table t rename to pharmacies;`); err != nil {
		level.Error(s.log).Log("method", "tx.ExecContext", "sql", "drop table pharmacies; alter table T rename to pharmacies;", "err", err)
		return errors.Wrap(ErrExecContextPharmaciesDB, err)
	}

	err := tx.Commit()
	if err != nil {
		level.Error(s.log).Log("method", "tx.Commit", "err", err)
		return errors.Wrap(ErrTxCommit, err)
	}

	return nil
}
