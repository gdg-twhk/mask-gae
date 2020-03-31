package postgres

import (
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/gomurphyx/sqlx"
	_ "github.com/lib/pq" // required for SQL access
	migrate "github.com/rubenv/sql-migrate"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

// Connect creates a connection to the PostgreSQL instance and applies any
// unapplied database migrations. A non-nil error is returned to indicate
// failure.
func Connect(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sqlx.Open("cloudsqlpostgres", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDB(db *sqlx.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "pharmacy_1",
				Up: []string{`
					create table if not exists pharmacies
					(
						id varchar(10) not null
							constraint pharmacies_pkey
								primary key,
						name varchar(254) not null,
						phone varchar(254) not null,
						address varchar(254) not null,
						mask_adult integer not null,
						mask_child integer not null,
						available varchar(1024) not null,
						note varchar(1024) not null,
						longitude double precision not null,
						latitude double precision not null,
						updated timestamp with time zone,
						custom_note varchar(1024) default ''::character varying not null,
						website varchar(1024) default ''::character varying not null
					);
					
					alter table pharmacies owner to postgres;

					CREATE EXTENSION if not exists cube;
					CREATE EXTENSION if not exists earthdistance;
					`,
				},
				Down: []string{"DROP TABLE pharmacies"},
			},
			{
				Id: "pharmacy_2",
				Up: []string{`
					alter table pharmacies
						add service_periods varchar(21) default '' not null;
					
					alter table pharmacies
						add service_note varchar(1024) default '' not null;
					
					alter table pharmacies
						add county varchar(19) default '' not null;
					
					alter table pharmacies
						add town varchar(10) default '' not null;
					
					alter table pharmacies
						add cunli varchar(10) default '' not null;
				`},
				Down: []string{`
					alter table pharmacies drop column service_periods;
					alter table pharmacies drop column service_note;
					alter table pharmacies drop column county;
					alter table pharmacies drop column town;
					alter table pharmacies drop column cunli;
				`},
			},
			{
				Id: "pharmacy_3",
				Up: []string{`
					alter table pharmacies alter column name set default '';
					alter table pharmacies alter column phone set default '';
					alter table pharmacies alter column address set default '';
					alter table pharmacies alter column mask_adult set default 0;
					alter table pharmacies alter column mask_child set default 0;
					alter table pharmacies alter column available set default '';
					alter table pharmacies alter column note set default '';
					alter table pharmacies alter column longitude set default 0.0;
					alter table pharmacies alter column latitude set default 0.0;
				`},
				Down: []string{``},
			},
			{
				Id: "pharmacy_4",
				Up: []string{`
					create or replace view latest_pharmacy_table as
					SELECT table_schema,
						   table_name
					FROM information_schema.tables
					WHERE table_type = 'BASE TABLE'
					  AND table_schema = 'public'
					  and table_name like 'pharmacy_%'
					order by table_name desc
					limit 1;

					CREATE OR REPLACE FUNCTION footgun(IN _tablename TEXT, IN _keepcount int)
						RETURNS void
						LANGUAGE plpgsql
					AS
					$$
					DECLARE
						row record;
					BEGIN
						FOR row IN
							SELECT table_schema,
								   table_name
							FROM information_schema.tables
							where table_name not in (
								SELECT table_name
								FROM information_schema.tables
								WHERE table_type = 'BASE TABLE'
								  AND table_schema = 'public'
								  and table_name like _tablename
								order by table_name desc
								limit _keepcount
							)
							  AND table_schema = 'public'
							  and table_name like _tablename
							LOOP
								EXECUTE 'DROP TABLE ' || quote_ident(row.table_schema) || '.' || quote_ident(row.table_name);
								RAISE INFO 'Dropped table: %', quote_ident(row.table_schema) || '.' || quote_ident(row.table_name);
							END LOOP;
					END;
					$$;
				`},
				Down: []string{},
			},
			{
				Id: "feedback_1",
				Up: []string{`
					create table if not exists options
					(
						id   varchar(21)  not null
							constraint options_pkey
								primary key,
						name varchar(254) not null
					);
					
					alter table options
						owner to postgres;

					INSERT INTO public.options (id, name) VALUES ('IRESxM58KC~dqg5XLCH~n', '自訂') ON CONFLICT (id) DO NOTHING;
					INSERT INTO public.options (id, name) VALUES ('ddCp1m88O4g5SU1GDJRPi', '當天已售完') ON CONFLICT (id) DO NOTHING;;
					INSERT INTO public.options (id, name) VALUES ('uYrYL~7Gd65IN2wWsWa9A', '號碼牌已發送完畢') ON CONFLICT (id) DO NOTHING;;
					INSERT INTO public.options (id, name) VALUES ('nAn6pj8UkrXST1syShrzV', '發放號碼牌') ON CONFLICT (id) DO NOTHING;;
				`},
				Down: []string{``},
			},
		},
	}

	_, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	return err
}
