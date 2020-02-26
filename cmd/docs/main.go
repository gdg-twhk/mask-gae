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
	"github.com/go-kit/kit/log/level"

	_ "github.com/cage1016/mask/cmd/docs/docs"
	"github.com/cage1016/mask/internal/app/docs/transports"
)

const (
	defServiceName string = "docs"
	defLogLevel    string = "error"
	defServiceHost string = "localhost"
	defHTTPPort    string = "8180"
	envServiceName string = "QS_DOCS_SERVICE_NAME"
	envLogLevel    string = "QS_DOCS_LOG_LEVEL"
	envServiceHost string = "QS_DOCS_SERVICE_HOST"
	envHTTPPort    string = "PORT"
)

type config struct {
	serviceName string
	logLevel    string
	serviceHost string
	httpPort    string
}

// Env reads specified environment variable. If no value has been found,
// fallback is returned.
func env(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// @title Mask API
// @version 0.2.0
// @description This is a Mask server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/tnstiger/mask-gdg/issues
// @contact.email cage.chung@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host mask.goodideas-studio.com
// @schemes https
// @BasePath /
func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = level.NewFilter(logger, level.AllowInfo())
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	cfg := loadConfig(logger)
	logger = log.With(logger, "service", cfg.serviceName)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//docs.SwaggerInfo.Host = cfg.serviceHost

	h := transports.NewHTTPHandler()

	wg := &sync.WaitGroup{}

	go startHTTPServer(ctx, wg, h, cfg.httpPort, logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	cancel()
	wg.Wait()

	fmt.Println("main: all goroutines have told us they've finished")
}

func loadConfig(logger log.Logger) (cfg config) {
	cfg.serviceName = env(envServiceName, defServiceName)
	cfg.logLevel = env(envLogLevel, defLogLevel)
	cfg.serviceHost = env(envServiceHost, defServiceHost)
	cfg.httpPort = env(envHTTPPort, defHTTPPort)
	return cfg
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, handler http.Handler, port string, logger log.Logger) {
	wg.Add(1)
	defer wg.Done()

	if port == "" {
		level.Error(logger).Log("protocol", "HTTP", "exposed", port, "err", "port is not assigned exist")
		return
	}

	p := fmt.Sprintf(":%s", port)
	// create a server
	srv := &http.Server{Addr: p, Handler: handler}
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
