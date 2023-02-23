package repositories

import (
	"PetProject/internal/useCases/models"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx"
)

type Balance struct {
	conn *pgx.Conn
}

func NewBalance(conn *pgx.Conn) *Balance {

	return &Balance{
		conn: conn,
	}
}

func (b *Balance) CreateUserBalance(m models.JSON) error {

	_, err := b.conn.Exec("INSERT INTO Balance (userID, balance) VALUES ($1::integer, $2::integer)", m.ID, m.Balance)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}
	return nil
}

func (b *Balance) CheckUserBalance(m models.JSON) bool {
	row := b.conn.QueryRow(
		"SELECT EXISTS (Select userID FROM balance Where userID = $1::integer);", m.ID)

	var res bool
	row.Scan(&res)

	return res
}

func (b *Balance) UpdateUserBalance(m models.JSON) error {

	_, err := b.conn.Exec("UPDATE Balance SET balance = balance + $1::integer WHERE userID = $2::integer;", m.Balance, m.ID)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}

	return nil
}

func (b *Balance) GetBalance(m models.JSON) string {

	row := b.conn.QueryRow(
		"Select sum(balance) FROM Balance Where userid = $1::integer", m.ID)

	var res int
	row.Scan(&res)

	fmt.Println(res)
	return strconv.Itoa(res)
}
