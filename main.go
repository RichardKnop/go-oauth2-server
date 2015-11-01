package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RichardKnop/go-microservice-example/service"
	// "github.com/RichardKnop/go-microservice-example/migrations"
	"github.com/gorilla/mux"
)

func main() {
	log.Print("Starting up the Go Microservice Example")

	// Run migrations
	// migrations.RunAll()

	r := mux.NewRouter()
	// Register /api/v1/tokens/ handler
	log.Print("Registering /api/v1/tokens/ resource handler")
	r.HandleFunc("/api/v1/tokens/", service.TokensHandler)

	// Listen on port 8080
	port := 8080
	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
