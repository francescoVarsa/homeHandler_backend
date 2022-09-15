package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func handleSecrets() {
	file, err := os.Open("./secrets.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()

		explodedRow := strings.Split(row, ":")
		key := explodedRow[0]
		secret := explodedRow[1]

		switch key {
		case "jwt_secret":
		case "DB_PASSWORD":
			os.Setenv(key, secret)
		}
	}

}
