package main

import (
	"log"
	"net/http"

	"realworld-go-nolambda/routes"
	"github.com/gorilla/mux"
)

func main() {


	// database.Setup()
	route := mux.NewRouter()
	routes.RegisterRoutes(route)

	if err := http.ListenAndServe(":8080", route); err != nil {
		log.Fatal(err)
	}

}
