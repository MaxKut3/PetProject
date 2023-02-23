package handlers

import (
	"PetProject/internal/useCases/models"
	"PetProject/internal/useCases/repositories"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx"
)

func HelloGoHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		Port:     5434,
		Database: "postgres",
		User:     "postgres",
		Password: os.Getenv("DB"),
	})
	if err != nil {
		log.Fatalf("Нет подключения к базе: %e", err)
	}

	Balance := repositories.NewBalance(conn)

	//Чтобы как-то проверить работу сделал (потом надо парсер JSON вставить)
	m := models.JSON{
		ID:      1,
		Balance: 1000,
	}

	fmt.Println(Balance.GetBalance(m))

	w.Write([]byte(Balance.GetBalance(m)))

}

func PutBalanceHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		Port:     5434,
		Database: "postgres",
		User:     "postgres",
		Password: os.Getenv("DB"),
	})
	if err != nil {
		log.Fatalf("Нет подключения к базе: %e", err)
	}

	Balance := repositories.NewBalance(conn)

	/*
		value, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m := models.JSON{}

		parseErr := models.Parser(value, m)
		if parseErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	*/
	m := models.JSON{
		ID:      2,
		Balance: 1000,
	}

	fmt.Println(Balance.CheckUserBalance(m))

	if !Balance.CheckUserBalance(m) {

		fmt.Println("Зашел в Create")
		Balance.CreateUserBalance(m)
		//		if createErr != nil {
		//			http.Error(w, err.Error(), http.StatusInternalServerError)
		//			return
		//		}
	} else {
		fmt.Println("Зашел в Update")
		Balance.UpdateUserBalance(m)
		//		if updateErr != nil {
		//			http.Error(w, err.Error(), http.StatusInternalServerError)
		//			return
		//		}
	}
	w.WriteHeader(http.StatusCreated)
}
