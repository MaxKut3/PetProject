package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (b *BalanceHandler) TransferMoneyHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		id1, _ := strconv.Atoi(chi.URLParam(r, "user1"))
		id2, _ := strconv.Atoi(chi.URLParam(r, "user2"))
		sum, _ := strconv.Atoi(chi.URLParam(r, "sum"))

		err := b.rep.TransferMoney(id1, id2, sum)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}
}
