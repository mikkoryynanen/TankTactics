package main

import (
	"fmt"
	"main/cmd/app"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	fmt.Println("starting server..")

	app := app.NewApp()
	app.Run()
}
