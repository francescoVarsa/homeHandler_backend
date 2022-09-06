package main

import (
	"github.com/gorilla/mux"
)

func (app *application) router() *mux.Router {
	router := mux.NewRouter()
	apiVersion := app.config.version

	router.HandleFunc("/"+apiVersion+"/getFoods", app.getFoodList).Methods("GET")

	return router
}
