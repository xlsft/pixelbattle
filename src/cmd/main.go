package main

import (
	"fmt"

	"github.com/joho/godotenv"

	"github.com/xlsft/pixelbattle/app"
	"github.com/xlsft/pixelbattle/database"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env file not found, skipping")
	}
	database.InitializeDatabase()
	app.StartService()
}
