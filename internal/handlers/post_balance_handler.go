package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (b *BalanceHandler) PostBalanceHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))
		balance, _ := strconv.Atoi(chi.URLParam(r, "balance"))

		err := b.rep.CreateUserBalance(userID, balance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}
