package main

import (
	"PetProject/config"
	"PetProject/internal/app"
)

func main() {

	cfg := config.NewConfig()
	app.Run(cfg)

}
