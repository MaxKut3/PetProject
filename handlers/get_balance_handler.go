package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MaxKut3/PetProject/repositories"

	"github.com/go-chi/chi/v5"
)

type BalanceResponse struct {
	UserID  int `json:"userID"`
	Balance int `json:"Balance"`
}

type BalanceHandler struct {
	rep *repositories.BalanceRepository
}

func NewBalanceHandler(rep *repositories.BalanceRepository) *BalanceHandler {
	return &BalanceHandler{
		rep: rep,
	}
}

func (b *BalanceHandler) GetBalanceHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, _ := strconv.Atoi(chi.URLParam(r, "userID"))

		ans := BalanceResponse{
			UserID:  userID,
			Balance: b.rep.GetBalance(userID),
		}

		resp, err := json.Marshal(&ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(resp)
		}
	}
}
