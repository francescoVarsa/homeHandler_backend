package main

import (
	"database/sql"
	"fmt"
	"homeHandler/models"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	mail "github.com/xhit/go-simple-mail/v2"
)

type config struct {
	version string
	port    int
}

type application struct {
	logger   *log.Logger
	config   config
	models   models.Models
	mailChan chan struct {
		To  string
		Msg string
	}
}

func main() {
	configureEnvironment()
	log.Println("opening db connection")
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
			port:    viper.GetInt("SERVER_PORT"),
		},
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Handler: app.router(),
		Addr:    fmt.Sprintf("0.0.0.0:%d", app.config.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Setup an email server
	server := mail.NewSMTPClient()
	// This is the IP address of the container that keeps the mailhog mail server used in
	// development. If the container is destroyed the IP will change and emails are not be
	// sended
	server.Host = viper.GetString("MAIL_SERVER_HOST")
	server.Port = viper.GetInt("MAIL_SERVER_PORT")
	server.Username = "admin@example.com"
	server.Password = "superSecretPassword"
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = true

	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	emailChannel := make(chan struct {
		To  string
		Msg string
	})
	app.mailChan = emailChannel

	go func() {
		for {
			msg := <-emailChannel
			SendMessage(msg.Msg, msg.To, smtpClient)
		}
	}()

	log.Printf("Server is running on port: %d", app.config.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func openDB() (*sql.DB, error) {
	dbConnString := fmt.Sprint(viper.Get("DB_STRING"))
	db, err := sql.Open("postgres", dbConnString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func configureEnvironment() {
	envType := os.Getenv("ENVIRONMENT")
	var envFile string

	if envType == "local" {
		envFile = ".env"
	} else if envType == "dev" {
		envFile = "dev.env"
	} else if envType == "prod" {
		envFile = "prod.env"
	}

	viper.SetConfigFile(envFile)
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("ttt")
		log.Fatal(err)
		return
	}

	log.Println("environment loaded")
}
