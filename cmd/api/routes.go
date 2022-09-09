package main

import (
	"github.com/gorilla/mux"
)

func (app *application) router() *mux.Router {
	router := mux.NewRouter()
	apiVersion := app.config.version

	router.HandleFunc("/"+apiVersion+"/users", app.getUsers).Methods("GET")

	return router
}
