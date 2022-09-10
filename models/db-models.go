package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) GetAllUsers() (*UsersList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select * from users`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var user User
	var users UsersList
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Password, &user.Email)

		if err != nil {
			return nil, err
		}

		nutriPlans, err := m.GetNutritionByUser(user.ID)

		if err != nil {
			log.Println(err)
			user.NutritionPlans = nil
		}

		user.NutritionPlans = *nutriPlans

		users = append(users, user)
	}

	return &users, nil
}

func (m *DBModel) GetNutritionByUser(userID int) (*[]NutritionPlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, plan_name from nutrition_plans where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var plans []NutritionPlan
	var plan NutritionPlan

	for rows.Next() {
		err := rows.Scan(&plan.ID, &plan.PlanName)

		if err != nil {
			return nil, err
		}

		foodList, err := m.GetPlanFood(plan.ID)

		if err != nil {
			return nil, err
		}

		plan.Foods = *foodList

		plans = append(plans, plan)
	}

	return &plans, nil
}

func (m *DBModel) GetPlanFood(planID int) (*FoodList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select food_name, meal_type, day_of_the_week from foodsList where plan_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, planID)

	if err != nil {
		return nil, err
	}

	var foodList FoodList
	var food Food

	for rows.Next() {
		err := rows.Scan(&food.Name, &food.MealType, &food.DayOfWeek)

		if err != nil {
			return nil, err
		}

		foodList = append(foodList, food)
	}

	return &foodList, nil
}

func (m *DBModel) AddFoodPlan(foodPlan *NutritionPlan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into nutrition_plans (user_id, plan_name, created_at, updated_at) values ($1, $2, $3, $4)`
	_, err := m.DB.ExecContext(ctx, query, foodPlan.UserID, foodPlan.PlanName, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) AddFood(food *Food) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into foodsList (food_name, plan_id, meal_type, day_of_the_week)
	 values
	  ($1, (select id from nutrition_plans where id = $2), $3, $4)`
	_, err := m.DB.ExecContext(ctx, query, food.Name, food.PlanID, food.MealType, food.DayOfWeek)

	if err != nil {
		return err
	}

	return nil
}
