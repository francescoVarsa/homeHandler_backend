package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	version string
	port    int
}

type application struct {
	logger *log.Logger
	config config
}

func main() {
	app := &application{
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		config: config{
			version: "v1",
			port:    4000,
		},
	}

	srv := &http.Server{
		Handler: app.router(),
		Addr:    fmt.Sprintf("127.0.0.1:%d", app.config.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is running on port: %d", app.config.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}
