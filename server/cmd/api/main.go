package main

import (
	"log"
	"os"

	"github.com/kaviraj-j/duoplay/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(application.Run())
}
