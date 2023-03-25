package handlers

import "net/http"

func (b *BalanceHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}
