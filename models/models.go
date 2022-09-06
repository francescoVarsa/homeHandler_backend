package models

type Food struct {
	ID             int    `json:"id"`
	FoodName       string `json:"food_name"`
	Quantity       int    `json:"quantity"`
	Calories       int    `json:"calories"`
	MacroNutrients string `json:"macro_nutrients"`
}

type FoodList []Food
