package app

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/MaxKut3/PetProject/internal/handlers"
	"github.com/MaxKut3/PetProject/internal/repositories"

	"github.com/MaxKut3/PetProject/config"

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

	balanceRep := repositories.NewBalanceRepository(conn)

	balanceHandler := handlers.NewBalanceHandler(balanceRep)

	r.Get("/", balanceHandler.HelloHandler)
	r.Get("/GetBalance/user/{userID}", balanceHandler.GetBalanceHandler())
	r.Get("/GetTransactionReport/user/{userID}", balanceHandler.GetTransactionsHandler())                //
	r.Get("/GetAccountingReport/month/{month}/year/{year}", balanceHandler.GetAccountingReportHandler()) //

	r.Post("/PostBalance/user/{userID}/balance/{balance}", balanceHandler.PostBalanceHandler())
	r.Post("/PostReserve/user/{userID}/amount/{amount}", balanceHandler.ResevreMoney())

	r.Put("/PutBalance/user/{userID}/balance/{balance}", balanceHandler.PutBalanceHandler())
	r.Get("/PutTransfer/user/{user1}/to/{user2}/sum/{sum}", balanceHandler.TransferMoneyHandler())
	r.Put("/PutRefresh/user/{userID}/", balanceHandler.RefreshMoneyHandler())

	r.Delete("/DeleteReserve/user/{userID}/", balanceHandler.DeleteReserveHandler())
	http.ListenAndServe(":8080", r)

}
