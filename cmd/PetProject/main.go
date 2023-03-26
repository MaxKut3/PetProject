package main

import (
	"github.com/MaxKut3/PetProject/config"
	"github.com/MaxKut3/PetProject/internal/app"
)

func main() {

	cfg := config.NewConfig()
	app.Run(cfg)

}
