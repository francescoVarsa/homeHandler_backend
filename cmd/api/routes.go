package main

import (
	"github.com/gorilla/mux"
)

func (app *application) router() *mux.Router {
	router := mux.NewRouter()
	apiVersion := app.config.version

	router.HandleFunc("/"+apiVersion+"/users", app.getUsers).Methods("GET")
	router.HandleFunc("/"+apiVersion+"/food/{id}", app.GetFood).Methods("GET")
	router.HandleFunc("/"+apiVersion+"/addPlan", app.NewFoodPlan).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/addFood", app.AddFoodToPlan).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/removePlan/{id}", app.RemoveFoodPlan).Methods("DELETE")
	router.HandleFunc("/"+apiVersion+"/updateFood/{id}", app.UpdateFood).Methods("PATCH")

	router.HandleFunc("/"+apiVersion+"/signUp", app.SignUp).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/signIn", app.SignIn).Methods("POST")

	return router
}
