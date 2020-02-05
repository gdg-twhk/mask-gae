package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/cage1016/mask/transport"
)

func main() {
	errs := make(chan error, 1)
	var err error

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	go startHTTPServer(transport.MakeHandler(), port, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	log.Fatal(fmt.Sprintf("http service terminated: %s", err))
}

func startHTTPServer(httpHandler http.Handler, port string, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	errs <- http.ListenAndServe(p, httpHandler)
}
