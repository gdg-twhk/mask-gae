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
	ErrQueryStoreFromPharmaciesDB   = errors.New("query pharmacies from DB failed")
	ErrInsertOrUpdateToPharmaciesDB = errors.New("insert or update DB failed")
	ErrExecContextPharmaciesDB      = errors.New("exec context pharmacies DB failed")
)

var _ model.PharmacyRepository = (*pharmacyRepository)(nil)

var t string

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
	if t == "" {
		buf, err := s.GetLatestPharmacyTableName(ctx)
		if err != nil {
			return []model.Pharmacy{}, err
		}
		t = buf
	}

	q := fmt.Sprintf(`SELECT *, point ($1, $2) <@> point(longitude, latitude)::point as distance
			FROM (select * from %s where longitude >= $3 and longitude <= $4 and latitude >= $5 and latitude <= $6) as a
			ORDER BY distance limit $7;`, t)

	pharmacies := []model.Pharmacy{}
	if err := s.db.SelectContext(ctx, &pharmacies, q, centerLng, centerLat, swLng, neLng, swLat, neLat, max); err != nil {
		level.Error(s.log).Log("method", "s.db.SelectContext", "err", err)
		return pharmacies, errors.Wrap(ErrQueryStoreFromPharmaciesDB, err)
	}
	return pharmacies, nil
}

func (s pharmacyRepository) Insert(ctx context.Context, updated string, pharmacies [][]model.Pharmacy) error {
	ctx = context.Background()
	nt := fmt.Sprintf("pharmacy_%s", updated)
	level.Info(s.log).Log("prepare", nt)

	var exists bool
	if err := s.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = $1);`, nt); err != nil {
		level.Error(s.log).Log("method", "s.db.GetContext", err, err)
		return err
	}

	if exists {
		level.Warn(s.log).Log("method", "s.db.GetContext", "table", nt, "msg", "exist and skip")
		return nil
	}

	q := fmt.Sprintf(`create table if not exists %s
								(
									id varchar(10),
									name varchar(254),
									phone varchar(254),
									address varchar(254),
									mask_adult integer,
									mask_child integer,
									available varchar(1024),
									note varchar(1024),
									longitude double precision,
									latitude double precision,
									updated timestamp with time zone,
									custom_note varchar(1024),
									website varchar(1024),
									service_periods varchar(21),
									service_note varchar(1024),
									county varchar(19),
									town varchar(10),
									cunli varchar(10)
								);
	
							alter table %s owner to postgres;`, nt, nt)
	if _, err := s.db.ExecContext(ctx, q); err != nil {
		level.Error(s.log).Log("method", "s.db.ExecContext", "sql", q, "err", err)
		return errors.Wrap(ErrExecContextPharmaciesDB, err)
	}

	q = fmt.Sprintf(`INSERT INTO %s (id, name, phone, address, mask_adult, mask_child, available, note, longitude, latitude, custom_note, website, updated, service_periods, service_note, county, town, cunli)
			VALUES (:id, :name, :phone, :address, :mask_adult, :mask_child, :available, :note, :longitude, :latitude, :custom_note, :website, :updated, :service_periods, :service_note, :county, :town, :cunli);`, nt)
	for _, trunk := range pharmacies {
		if _, err := s.db.NamedExecContext(ctx, q, trunk); err != nil {
			level.Error(s.log).Log("method", "s.db.NamedExecContext", "err", err)
			return errors.Wrap(ErrInsertOrUpdateToPharmaciesDB, err)
		}
	}

	// update the latest table_name to memory
	buf, err := s.GetLatestPharmacyTableName(ctx)
	if err != nil {
		return err
	}
	t = buf

	return nil
}

func (s pharmacyRepository) GetLatestPharmacyTableName(ctx context.Context) (string, error) {
	lt := struct {
		TableName string `db:"table_name"`
	}{}

	q := `select table_name from latest_pharmacy_table;`
	if err := s.db.GetContext(ctx, &lt, q); err != nil {
		return "", err
	}
	return lt.TableName, nil
}

func (s pharmacyRepository) FootGun(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `select footgun('pharmacy_%', 5)`); err != nil {
		level.Error(s.log).Log("method", "s.db.ExecContext", "sql", `select footgun('pharmacy_%', 5)`, "err", err)
		return err
	}
	return nil
}
