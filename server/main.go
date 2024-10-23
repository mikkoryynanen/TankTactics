package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	fmt.Println("starting server..")
	app := NewApp()
	// var app App
	app.Run()
}
