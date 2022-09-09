package models

import "database/sql"

//Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModels returns models with db pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

type Food struct {
	ID             int    `json:"id"`
	FoodName       string `json:"food_name"`
	Quantity       int    `json:"quantity"`
	Calories       int    `json:"calories"`
	MacroNutrients string `json:"macro_nutrients"`
	MealType       string `json:"meal_type"`
}

type FoodList []Food

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UsersList []User
