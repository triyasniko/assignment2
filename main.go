package main

import (
	"assignment-2/config"
	"assignment-2/controllers"
	"assignment-2/routers"
)

var PORT = ":8080"

func main() {
	db := config.StartDB()
	controller := controllers.New(db)

	routers.StartServer(controller).Run(PORT)
}
