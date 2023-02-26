package app

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MaxKut3/PetProject/config"

	"github.com/MaxKut3/PetProject/internal/use_cases/handlers"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run(cfg *config.Config) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env файл не прочитан: %e", err)
	}

	port, _ := strconv.ParseUint(os.Getenv("DBPORT"), 10, 16)

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     os.Getenv("DBHOST"),
		Port:     uint16(port),
		Database: os.Getenv("DB"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASSWORD"),
	})
	if err != nil {
		log.Fatalf("Нет подключения к базе: %e", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	balanceRep := handlers.NewBalanceRepository(conn)

	helloHandler := handlers.NewHelloHandler()
	r.Get("/", helloHandler)

	postHandler := handlers.NewPostBalanceHandler(balanceRep)
	r.Post("/PostBalance/user/{userID}/balance/{balance}", postHandler)

	putBalanceHandler := handlers.NewPutBalanceHandler(balanceRep)
	r.Put("/PutBalance/user/{userID}/balance/{balance}", putBalanceHandler)

	getBalanceHandler := handlers.NewGetBalanceHandler(balanceRep)
	r.Get("/GetBalance/user/{userID}", getBalanceHandler)

	http.ListenAndServe(":8080", r)

}
