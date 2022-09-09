package main

import (
	"context"
	"encoding/json"
	"homeHandler/models"
	"net/http"
	"time"
)

func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	query := `select * from users`

	rows, err := app.models.DB.DB.QueryContext(ctx, query)

	if err != nil {
		app.logger.Println(err)
		return
	}

	var user models.User
	var users []models.User
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.Password, &user.Email)

		if err != nil {
			app.logger.Println(err)
			return
		}

		users = append(users, user)
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
