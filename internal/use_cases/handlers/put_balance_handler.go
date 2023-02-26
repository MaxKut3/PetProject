package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx"
)

func (rep *BalanceRepository) FindUserBalance(userID int) bool {

	row := rep.conn.QueryRow(
		"SELECT EXISTS (Select userID FROM balance Where userID = $1::integer);", userID)

	var res bool
	row.Scan(&res)

	return res
}

func (rep *BalanceRepository) UpdateUserBalance(userID, balance int) error {

	_, err := rep.conn.Exec("UPDATE Balance SET balance = balance + $1::integer WHERE userID = $2::integer;", balance, userID)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}
	return nil
}

func (rep *BalanceRepository) AddTransaction(userID, balance int) error {

	_, err := rep.conn.Exec("INSERT INTO Transactions (type_transaction, userID, amount) VALUES (0, $1, $2);", userID, balance)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}
	return nil
}

type putBalanceHandler struct {
	conn *pgx.Conn
}

func NewPutBalanceHandler(rep *BalanceRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		id, _ := strconv.Atoi(chi.URLParam(r, "userID"))
		balance, _ := strconv.Atoi(chi.URLParam(r, "balance"))

		errBalance := rep.UpdateUserBalance(id, balance)
		errTransaction := rep.AddTransaction(id, balance)
		if errBalance != nil || errTransaction != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}
}
