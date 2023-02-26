package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx"
)

func (rep *BalanceRepository) CreateUserBalance(userID, balance int) error {

	_, err := rep.conn.Exec("INSERT INTO Balance (userID, balance) VALUES ($1::integer, $2::integer)", userID, balance)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}
	return nil
}

type postBalanceHandler struct {
	conn *pgx.Conn
}

func NewPostBalanceHandler(rep *BalanceRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		id, _ := strconv.Atoi(chi.URLParam(r, "userID"))
		balance, _ := strconv.Atoi(chi.URLParam(r, "balance"))

		errBalance := rep.CreateUserBalance(id, balance)
		if errBalance != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			rep.AddTransaction(id, balance)
			w.WriteHeader(http.StatusCreated)
		}
	}
}
