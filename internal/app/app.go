package app

import (
	"PetProject/config"
	"PetProject/internal/useCases/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(cfg *config.Config) {

	//Непонял как закидывать conn в handler ?

	/*
		conn, err := pgx.Connect(cfg.Conn)
		if err != nil {
			log.Fatalf("Нет подключения к базе: %e", err)
		}

		Balance := repositories.NewBalance(conn)

	*/

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.HandleFunc("/", handlers.HelloGoHandler)
	r.HandleFunc("/GetBalance", handlers.GetBalanceHandler)
	r.HandleFunc("/PutBalance", handlers.PutBalanceHandler)

	http.ListenAndServe(":8080", r)

}
