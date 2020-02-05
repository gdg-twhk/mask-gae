package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {	
	fmt.Fprint(w, "ownershiop")
}

func webmasterHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./googleae8f4bcce8bec00c.html")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/googleae8f4bcce8bec00c.html", webmasterHandler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	// [END setting_port]
}
