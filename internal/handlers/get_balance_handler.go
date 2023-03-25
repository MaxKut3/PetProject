package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MaxKut3/PetProject/internal/repositories"

	"github.com/go-chi/chi/v5"
)

type BalanceResponse struct {
	User_id int `json:"user_id"`
	Balance int `json:"balance"`
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
			User_id: userID,
			Balance: b.rep.GetBalance(userID),
		}

		resp, err := json.Marshal(&ans)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(resp)

	}
}
