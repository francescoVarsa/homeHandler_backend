package main

import (
	"database/sql"
	"errors"
	"fmt"
	"homeHandler/models"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	version string
	port    int
}

type application struct {
	logger *log.Logger
	config config
	models models.Models
}

func main() {
	db, err := openDB()

	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		config: config{
			version: "v1",
			port:    4000,
		},
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Handler: app.router(),
		Addr:    fmt.Sprintf("127.0.0.1:%d", app.config.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is running on port: %d", app.config.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

func openDB() (*sql.DB, error) {
	pwd := os.Getenv("DB_PASSWORD")

	if pwd == "" {
		return nil, errors.New("no password for the database was provided")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("postgres://francesco:%s@localhost/homeManager?sslmode=disable", pwd))

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
