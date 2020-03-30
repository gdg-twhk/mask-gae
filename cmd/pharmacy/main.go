package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gomurphyx/sqlx"

	"github.com/cage1016/mask/internal/app/pharmacy/endpoints"
	"github.com/cage1016/mask/internal/app/pharmacy/postgres"
	"github.com/cage1016/mask/internal/app/pharmacy/service"
	"github.com/cage1016/mask/internal/app/pharmacy/transports"
	"github.com/cage1016/mask/internal/pkg/level"
	psql "github.com/cage1016/mask/internal/pkg/postgres"
)

const (
	defServiceName   = "pharmacy"
	defLogLevel      = "error"
	defServiceHost   = "localhost"
	defHTTPPort      = "8080"
	defDBHost        = ""
	defDBPort        = ""
	defDBUser        = ""
	defDBPass        = ""
	defDBName        = ""
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""

	envServiceName   = "SERVICE_NAME"
	envLogLevel      = "LOG_LEVEL"
	envServiceHost   = "SERVICE_HOST"
	envHTTPPort      = "PORT"
	envDBHost        = "DB_HOST"
	envDBPort        = "DB_PORT"
	envDBUser        = "DB_USER"
	envDBPass        = "DB_PASS"
	envDBName        = "DB"
	envDBSSLMode     = "DB_SSL_MODE"
	envDBSSLCert     = "DB_SSL_CERT"
	envDBSSLKey      = "DB_SSL_KEY"
	envDBSSLRootCert = "DB_SSL_ROOT_CERT"
)

type config struct {
	serviceName string
	logLevel    string
	serviceHost string
	httpPort    string
	dbConfig    psql.Config
}

// Env reads specified environment variable. If no value has been found,
// fallback is returned.
func env(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowWarn())
		logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	cfg := loadConfig(logger)
	logger = log.With(logger, "service", cfg.serviceName)
	level.Info(logger).Log("version", service.Version, "commitHash", service.CommitHash, "buildTimeStamp", service.BuildTimeStamp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	svc := NewServer(db, logger)
	eps := endpoints.New(svc, logger)

	wg := &sync.WaitGroup{}

	go startHTTPServer(ctx, wg, eps, cfg.httpPort, logger)
	go tickerFunc(ctx, wg, svc, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	cancel()
	wg.Wait()

	fmt.Println("main: all goroutines have told us they've finished")
}

func loadConfig(_ log.Logger) (cfg config) {
	dbConfig := psql.Config{
		Host:        env(envDBHost, defDBHost),
		Port:        env(envDBPort, defDBPort),
		User:        env(envDBUser, defDBUser),
		Pass:        env(envDBPass, defDBPass),
		Name:        env(envDBName, defDBName),
		SSLMode:     env(envDBSSLMode, defDBSSLMode),
		SSLCert:     env(envDBSSLCert, defDBSSLCert),
		SSLKey:      env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: env(envDBSSLRootCert, defDBSSLRootCert),
	}

	cfg.dbConfig = dbConfig
	cfg.serviceName = env(envServiceName, defServiceName)
	cfg.logLevel = env(envLogLevel, defLogLevel)
	cfg.serviceHost = env(envServiceHost, defServiceHost)
	cfg.httpPort = env(envHTTPPort, defHTTPPort)
	return cfg
}

func connectToDB(cfg psql.Config, logger log.Logger) *sqlx.DB {
	db, err := psql.Connect(cfg)
	if err != nil {
		level.Error(logger).Log(
			"host", cfg.Host,
			"port", cfg.Port,
			"user", cfg.User,
			"dbname", cfg.Name,
			"password", cfg.Pass,
			"sslmode", cfg.SSLMode,
			"SSLCert", cfg.SSLCert,
			"SSLKey", cfg.SSLKey,
			"SSLRootCert", cfg.SSLRootCert,
			"err", err,
		)
		os.Exit(1)
	}
	return db
}

func NewServer(db *sqlx.DB, logger log.Logger) service.PharmacyService {
	repo := postgres.New(db, logger)
	return service.New(repo, logger)
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, endpoints endpoints.Endpoints, port string, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	if port == "" {
		level.Error(logger).Log("protocol", "HTTP", "exposed", port, "err", "port is not assigned exist")
		return
	}

	p := fmt.Sprintf(":%s", port)
	// create a server
	srv := &http.Server{Addr: p, Handler: transports.NewHTTPHandler(endpoints, logger)}
	level.Info(logger).Log("protocol", "HTTP", "exposed", port)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			level.Info(logger).Log("Listen", err)
		}
	}()

	<-ctx.Done()

	// shut down gracefully, but wait no longer than 5 seconds before halting
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ignore error since it will be "Err shutting down server : context canceled"
	srv.Shutdown(shutdownCtx)

	level.Info(logger).Log("protocol", "HTTP", "Shutdown", "http server gracefully stopped")
}

func tickerFunc(ctx context.Context, wg *sync.WaitGroup, svc service.PharmacyService, logger log.Logger) {
	wg.Add(1)
	// tell the caller we've stopped
	defer wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			svc.TickerUpdate(ctx)
		case <-ctx.Done():
			logger.Log("method", "ctx.Done")
			return
		}
	}
}
