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

	query := `select food_name from foodsList where plan_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, planID)

	if err != nil {
		return nil, err
	}

	var foodList FoodList
	var foodName Food

	for rows.Next() {
		err := rows.Scan(&foodName.Name)

		if err != nil {
			return nil, err
		}

		foodList = append(foodList, foodName)
	}

	return &foodList, nil
}
