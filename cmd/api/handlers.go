package main

import (
	"encoding/json"
	"homeHandler/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	fPlan.PlanName = payload.PlanName

	err = app.models.DB.AddFoodPlan(&fPlan)

	if err != nil {
		app.logger.Println(err)
		return
	}

	jsonRes, err := json.MarshalIndent(&fPlan, "", " ")

	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonRes)

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
		app.logger.Println("=>", err)
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
			app.logger.Println(err)
		}

		foodEntries = append(foodEntries, foodEntry)
	}

	if err != nil {
		return
	}

	res, err := json.MarshalIndent(foodEntries, "", " ")

	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (app *application) RemoveFoodPlan(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	plan_id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.logger.Println(err)
		return
	}

	// foodList, err := app.models.DB.GetPlanFood(plan_id)

	if err != nil {
		app.logger.Println(err)
		return
	}

	// delete all foods related to the plan
	err = app.models.DB.RemoveFood(plan_id)

	if err != nil {
		app.logger.Println(err)
		return
	}

	err = app.models.DB.DeleteNutritionPlan(plan_id)

	if err != nil {
		app.logger.Println(err)
		return
	}

	var res struct {
		OK      bool
		message string
	}

	res.OK = true
	res.message = ""

	response, err := json.MarshalIndent(&res, "", " ")

	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (app *application) UpdateFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.logger.Println(err)
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
		app.logger.Println(err)
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
				app.logger.Println(err)
				return
			}
		}
	}

	updatedFood, err := app.models.DB.GetFoodByID(id)

	if err != nil {
		app.logger.Println(err)
		return
	}

	res, err := json.MarshalIndent(&updatedFood, "", " ")

	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (app *application) GetFood(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		app.logger.Println(err)
		return
	}

	food, err := app.models.DB.GetFoodByID(id)

	if err != nil {
		app.logger.Println(err)
		return
	}

	res, err := json.MarshalIndent(&food, "", " ")

	if err != nil {
		app.logger.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
