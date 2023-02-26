package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/jackc/pgx"
)

type getBalanceResponse struct {
	userID  int `json:"userID"`
	balance int `json:"Balance"`
}

type BalanceRepository struct {
	conn *pgx.Conn
}

func NewBalanceRepository(conn *pgx.Conn) *BalanceRepository {
	return &BalanceRepository{
		conn: conn,
	}
}

func (r *BalanceRepository) GetBalance(id int) int {
	row := r.conn.QueryRow(
		"Select sum(balance) FROM Balance Where userid = $1::integer", id)

	var balance int
	row.Scan(&balance)

	return balance
}

type GetBalanceHandler struct {
	UserID  int
	Balance int
	rep     *BalanceRepository
}

func NewGetBalanceHandler(rep *BalanceRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		id, _ := strconv.Atoi(chi.URLParam(r, "userID"))

		ans := getBalanceResponse{
			userID:  id,
			balance: rep.GetBalance(id),
		}

		fmt.Println(ans)

		resp, err := json.Marshal(&ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(resp)
		}
		fmt.Println(resp)
	}
}
