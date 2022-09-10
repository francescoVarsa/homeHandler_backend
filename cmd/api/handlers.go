package main

import (
	"encoding/json"
	"homeHandler/models"
	"log"
	"net/http"
)

func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.DB.GetAllUsers()

	if err != nil {
		//Handle error in a better way
		log.Println(err)
		return
	}

	usersJson, err := json.MarshalIndent(users, "", " ")

	if err != nil {
		app.logger.Println(err)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(usersJson)
}

type FoodPlanPayload struct {
	UserID   int           `json:"user_id"`
	PlanName string        `json:"plan_name"`
	Foods    []models.Food `json:"foods"`
}

func (app *application) NewFoodPlan(w http.ResponseWriter, r *http.Request) {
	var payload FoodPlanPayload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		log.Println(err)
		return
	}

	var fPlan models.NutritionPlan

	fPlan.UserID = payload.UserID
	fPlan.Foods = payload.Foods
	fPlan.PlanName = payload.PlanName

	_, err = app.models.DB.AddFoodPlan(&fPlan)

	if err != nil {
		app.logger.Println(err)
		return
	}

	type response struct {
		OK      bool
		message string
	}

	var res response

	res.OK = true
	res.message = ""

	jsonRes, err := json.MarshalIndent(&res, "", " ")

	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonRes)

}
