package main

import (
	"encoding/json"
	"log"
	"net/http"
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
