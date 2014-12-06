package main

import (
	"codemo"
)

func main() {
	app, _ := ideat.NewApp()
	// Listen and server on 0.0.0.0:8080
	app.Run(":8080")
}
