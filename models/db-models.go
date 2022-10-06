package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) GetAllUsers() (*UsersList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, last_name, email from users`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	var user User
	var users UsersList
	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.LastName,
			&user.Email)

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

func (m *DBModel) DeleteNutritionPlan(plan_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from nutrition_plans where id = $1`

	_, err := m.DB.ExecContext(ctx, query, plan_id)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) RemoveFood(plan_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from foodslist where plan_id = $1`

	_, err := m.DB.ExecContext(ctx, query, plan_id)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetPlanByID(planID int) (*NutritionPlan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, user_id, plan_name, created_at, updated_at from nutrition_plans where id = $1`
	row := m.DB.QueryRowContext(ctx, query, planID)

	var plan NutritionPlan

	err := row.Scan(&plan.ID, &plan.UserID, &plan.PlanName, &plan.CreatedAt, &plan.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &plan, nil

}

func (m *DBModel) UpdateFood(colName string, value string, foodID int) (*Food, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("update foodsList set %s = $1 where id = $2", colName)

	_, err := m.DB.ExecContext(ctx, query, value, foodID)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *DBModel) GetFoodByID(id int) (*Food, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select plan_id, food_name, meal_type, day_of_the_week from foodslist where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var food Food

	err := row.Scan(&food.PlanID, &food.Name, &food.MealType, &food.DayOfWeek)

	if err != nil {
		return nil, err
	}

	return &food, nil
}

func (m *DBModel) CreateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into users (name, last_name, password, email) values ($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, query, user.Name, user.LastName, user.Password, user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetUserPassword(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select password from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var pwd string

	err := row.Scan(&pwd)

	if err != nil {
		return "", err
	}

	return pwd, nil
}

func (m *DBModel) CheckExistingUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select email from users where email = $1`
	row := m.DB.QueryRowContext(ctx, query, username)

	var userEmail string
	err := row.Scan(&userEmail)

	if err != nil || len(userEmail) == 0 {
		return false
	}

	return true

}

func (m *DBModel) GetUserByUsername(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select email, last_name, id, name from users where email = $1`
	row := m.DB.QueryRowContext(ctx, query, username)

	var user User
	err := row.Scan(&user.Email, &user.LastName, &user.ID, &user.Name)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *DBModel) SetResetToken(id int, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set reset_token = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, token, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetUserResetToken(id int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select reset_token from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var resetToken string

	err := row.Scan(&resetToken)

	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (m *DBModel) SetNewPassword(id int, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set reset_token = $1, password = $2, updated_at = $3 where id = $4`

	_, err := m.DB.ExecContext(ctx, query, nil, newPassword, time.Now(), id)

	if err != nil {
		return err
	}

	return nil
}
