package main

import (
	"github.com/xlsft/pixelbattle/app"
	"github.com/xlsft/pixelbattle/database"
)

func main() {
	database.InitializeDatabase()
	app.StartService()
}
