package main

import (
	"PetProject/config"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.NewConfig()
	fmt.Println(cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	http.ListenAndServe(":8080", r)

}
