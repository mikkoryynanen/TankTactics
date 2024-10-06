package main

import (
	"fmt"
)

func main() {
	fmt.Println("starting server..")
	app := NewApp()
	// var app App
	app.Run()
}
