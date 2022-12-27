package main

import (
	"encoding/json"
	"homeHandler/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.DB.GetAllUsers()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, users, "data")
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
		app.errorJSON(w, err)
		return
	}

	var fPlan models.NutritionPlan

	fPlan.UserID = payload.UserID
	fPlan.PlanName = payload.PlanName

	err = app.models.DB.AddFoodPlan(&fPlan)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	type responseOK struct {
		OK      bool
		Message string
	}

	app.writeJSON(w, http.StatusOK, responseOK{
		OK:      true,
		Message: "Plan correctly created",
	}, "data")
}

func (app *application) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user, err := app.models.DB.GetUserByID(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, &user, "data")
}

type FoodPayload struct {
	Foods []struct {
		PlanID    int    `json:"plan_id"`
		Name      string `json:"food_name"`
		MealType  string `json:"meal_type"`
		DayOfWeek string `json:"day_of_the_week"`
	}
}

func (app *application) AddFoodToPlan(w http.ResponseWriter, r *http.Request) {
	var payload FoodPayload
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	var foodEntry models.Food
	var foodEntries models.FoodList

	for _, f := range payload.Foods {
		foodEntry.PlanID = f.PlanID
		foodEntry.Name = f.Name
		foodEntry.MealType = f.MealType
		foodEntry.DayOfWeek = f.DayOfWeek

		err = app.models.DB.AddFood(&foodEntry)

		if err != nil {
			app.errorJSON(w, err)
		}

		foodEntries = append(foodEntries, foodEntry)
	}

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, foodEntries, "data")

}

func (app *application) RemoveFoodPlan(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	plan_id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// delete all foods related to the plan
	err = app.models.DB.RemoveFood(plan_id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteNutritionPlan(plan_id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusNoContent, nil, "data")
}

func (app *application) UpdateFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	type reqBody struct {
		FoodName  string `json:"food_name"`
		DayOfWeek string `json:"day_of_the_week"`
		MealType  string `json:"meal_type"`
	}

	var payload reqBody

	err = json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	newFood := make(map[string]string)

	// Only those are the fields that user can update
	newFood["food_name"] = payload.FoodName
	newFood["day_of_the_week"] = payload.DayOfWeek
	newFood["meal_type"] = payload.MealType

	for key, val := range newFood {
		if len(val) != 0 {
			_, err := app.models.DB.UpdateFood(key, val, id)

			if err != nil {
				app.errorJSON(w, err)
				return
			}
		}
	}

	updatedFood, err := app.models.DB.GetFoodByID(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedFood, "data")
}

func (app *application) GetFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	food, err := app.models.DB.GetFoodByID(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, food, "data")

}
