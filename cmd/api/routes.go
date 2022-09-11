package main

import (
	"github.com/gorilla/mux"
)

func (app *application) router() *mux.Router {
	router := mux.NewRouter()
	apiVersion := app.config.version

	router.HandleFunc("/"+apiVersion+"/users", app.getUsers).Methods("GET")
	router.HandleFunc("/"+apiVersion+"/addPlan", app.NewFoodPlan).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/addFood", app.AddFoodToPlan).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/removePlan/{id}", app.RemoveFoodPlan).Methods("DELETE")

	return router
}
