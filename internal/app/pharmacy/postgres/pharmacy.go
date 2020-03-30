package postgres

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gomurphyx/sqlx"

	"github.com/cage1016/mask/internal/app/pharmacy/model"
	"github.com/cage1016/mask/internal/pkg/errors"
	"github.com/cage1016/mask/internal/pkg/level"
)

var (
	ErrQueryStoreFromPharmaciesDB = errors.New("query pharmacies from DB failed")
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

func (s pharmacyRepository) Query(ctx context.Context, latestPharmacyTable string, centerLng, centerLat, swLng, neLng, swLat, neLat float64, max uint64) ([]model.Pharmacy, error) {
	q := fmt.Sprintf(`SELECT *, point ($1, $2) <@> point(longitude, latitude)::point as distance
			FROM (select * from %s where longitude >= $3 and longitude <= $4 and latitude >= $5 and latitude <= $6) as a
			ORDER BY distance limit $7;`, latestPharmacyTable)

	pharmacies := []model.Pharmacy{}
	if err := s.db.SelectContext(ctx, &pharmacies, q, centerLng, centerLat, swLng, neLng, swLat, neLat, max); err != nil {
		level.Error(s.log).Log("method", "s.db.SelectContext", "err", err)
		return pharmacies, errors.Wrap(ErrQueryStoreFromPharmaciesDB, err)
	}
	return pharmacies, nil
}

func (s pharmacyRepository) GetLatestPharmacyTableName(ctx context.Context) (string, error) {
	lt := struct {
		TableName string `db:"table_name"`
	}{}

	q := `select table_name from latest_pharmacy_table;`
	if err := s.db.GetContext(ctx, &lt, q); err != nil {
		level.Error(s.log).Log("method", "GetLatestPharmacyTableName", "err", err)
		return "", err
	}
	return lt.TableName, nil
}
