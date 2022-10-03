package models

import (
	"database/sql"
	"time"
)

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

type FoodList []Food

type User struct {
	ID             int             `json:"id"`
	Name           string          `json:"name"`
	LastName       string          `json:"last_name"`
	Password       string          `json:"password"`
	Email          string          `json:"email"`
	NutritionPlans []NutritionPlan `json:"nutrition_plans"`
}

type NutritionPlan struct {
	ID        int       `json:"plan_id"`
	UserID    int       `json:"-"`
	PlanName  string    `json:"plan_name"`
	Foods     []Food    `json:"foods"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Food struct {
	PlanID    int    `json:"-"`
	Name      string `json:"food_name"`
	MealType  string `json:"meal_type"`
	DayOfWeek string `json:"day_of_the_week"`
}

type UsersList []User

type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
