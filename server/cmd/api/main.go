package main

import (
	"log"

	"github.com/kaviraj-j/duoplay/internal/app"
	"github.com/kaviraj-j/duoplay/internal/config"
)

func main() {
	// load config
	config, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	// create new app with config
	app, err := app.Create(&config)
	if err != nil {
		log.Fatal(err)
	}

	// run app
	app.Run()
}
