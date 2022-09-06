package main

import (
	"encoding/json"
	"homeHandler/models"
	"log"
	"net/http"
)

func (app *application) getFoodList(w http.ResponseWriter, r *http.Request) {
	var apple models.Food
	var pasta models.Food

	apple.ID = 1
	apple.FoodName = "Apple"
	apple.Quantity = 2
	apple.MacroNutrients = "sugar"
	apple.Calories = 230

	pasta.ID = 2
	pasta.FoodName = "Pasta"
	pasta.Quantity = 150
	pasta.MacroNutrients = "carbohydrates"
	pasta.Calories = 300

	var foodList = models.FoodList{apple, pasta}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.MarshalIndent(foodList, "", " ")

	if err != nil {
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
