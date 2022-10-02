package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) wrapMiddleware(next http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	}
}

func (app *application) router() http.Handler {
	router := mux.NewRouter()
	apiVersion := app.config.version
	secure := alice.New(app.checkToken)

	router.Use(app.enableCORS)

	// Secured routes
	router.HandleFunc("/"+apiVersion+"/users", app.wrapMiddleware(secure.ThenFunc(app.getUsers))).Methods("GET")
	router.HandleFunc("/"+apiVersion+"/food/{id}", app.wrapMiddleware(secure.ThenFunc(app.GetFood))).Methods("GET")
	router.HandleFunc("/"+apiVersion+"/addPlan", app.wrapMiddleware(secure.ThenFunc(app.NewFoodPlan))).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/addFood", app.wrapMiddleware(secure.ThenFunc(app.AddFoodToPlan))).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/removePlan/{id}", app.wrapMiddleware(secure.ThenFunc(app.RemoveFoodPlan))).Methods("DELETE")
	router.HandleFunc("/"+apiVersion+"/updateFood/{id}", app.wrapMiddleware(secure.ThenFunc(app.UpdateFood))).Methods("PATCH")

	router.HandleFunc("/"+apiVersion+"/signUp", app.SignUp).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/signIn", app.SignIn).Methods("POST")
	router.HandleFunc("/"+apiVersion+"/resetPassword", app.resetPassword).Methods("POST")

	return router
}
