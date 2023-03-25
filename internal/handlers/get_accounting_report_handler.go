package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (b *BalanceHandler) GetAccountingReportHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		month, _ := strconv.Atoi(chi.URLParam(r, "month"))
		year, _ := strconv.Atoi(chi.URLParam(r, "year"))

		err := b.rep.GetAccountingReport(month, year)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}
}
