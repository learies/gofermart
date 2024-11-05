package main

import (
	"log"
	"net/http"

	"github.com/learies/gofermart/internal/app"
)

func main() {
	application := app.NewApp()
	defer application.Close()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", application.Routes); err != nil {
		log.Fatal(err)
	}
}
