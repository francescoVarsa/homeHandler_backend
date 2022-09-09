package models

import (
	"context"
	"database/sql"
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

		users = append(users, user)
	}

	return &users, nil
}
