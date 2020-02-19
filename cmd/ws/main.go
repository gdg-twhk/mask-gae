package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	broker "github.com/nats-io/nats.go"

	adapter "github.com/cage1016/mask/internal/app/ws"
	"github.com/cage1016/mask/internal/app/ws/api"
	"github.com/cage1016/mask/internal/app/ws/nats"
	"github.com/cage1016/mask/internal/pkg/level"
)

const (
	defServiceName = "ws"
	defLogLevel    = "error"
	defPort        = "8180"
	defNatsURL     = broker.DefaultURL

	envServiceName = "QS_WS_SERVICE_NAME"
	envLogLevel    = "QS_WS_LOG_LEVEL"
	envPort        = "PORT"
	envNatsURL     = "QS_NATS_URL"
)

type config struct {
	serviceName string
	logLevel    string
	port        string
	natsURL     string
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
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	cfg := loadConfig(logger)

	logger = log.With(logger, "service", "ws")

	nc, err := broker.Connect(cfg.natsURL)
	if err != nil {
		level.Error(logger).Log("natsURL", cfg.natsURL, "err", err)
		os.Exit(1)
	}
	defer nc.Close()

	pubsub := nats.New(nc, logger)
	svc := newService(pubsub, logger)

	errs := make(chan error, 2)

	go func() {
		p := fmt.Sprintf(":%s", cfg.port)
		level.Info(logger).Log("port", cfg.port)
		errs <- http.ListenAndServe(p, api.MakeHandler(svc, logger))
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	level.Error(logger).Log("terminated", err)
}

func loadConfig(logger log.Logger) (cfg config) {
	cfg.serviceName = env(envServiceName, defServiceName)
	cfg.logLevel = env(envLogLevel, defLogLevel)
	cfg.port = env(envPort, defPort)
	cfg.natsURL = env(envNatsURL, defNatsURL)

	return cfg
}

func newService(pubsub adapter.Service, logger log.Logger) adapter.Service {
	svc := adapter.New(pubsub)
	svc = api.LoggingMiddleware(svc, logger)

	return svc
}
