package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on port %s", port)

	// Fire up server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
