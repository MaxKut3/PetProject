package repositories

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx"
)

type BalanceRepository struct {
	conn *pgx.Conn
}

func NewBalanceRepository(conn *pgx.Conn) *BalanceRepository {
	return &BalanceRepository{
		conn: conn,
	}
}

// addTransaction - Запись транзакции - должна быть во всех операциях, связанных с изменениями
// type_transaction - 0 пополнение, 1 - списание
func (rep *BalanceRepository) addTransaction(userID, balance, typeTrans int) error {

	_, err := rep.conn.Exec("INSERT INTO Transactions (userID, amount, type_transaction) VALUES ($1, $2, $3);", userID, balance, typeTrans)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}
	return nil
}

// balanceCheck Проверка достаточности средств
func (rep *BalanceRepository) balanceCheck(userID, amount int) bool {

	row := rep.conn.QueryRow(
		"Select sum(balance) FROM Balance Where userid = $1::integer", userID)

	var balance int
	row.Scan(&balance)

	if amount > balance {
		return false
	}
	return true
}

// CreateUserBalance - Создание баланса

func (rep *BalanceRepository) CreateUserBalance(userID, balance int) error {

	tx, _ := rep.conn.Begin()

	errTrans := rep.addTransaction(userID, balance, 0)
	_, errInsert := rep.conn.Exec("INSERT INTO Balance (userID, balance) VALUES ($1::integer, $2::integer)", userID, balance)

	if errTrans != nil || errInsert != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка записи в базу")
	}

	tx.Commit()
	return nil

}

func (rep *BalanceRepository) updateUserBalance(userID, balance int) error {

	_, errUpdate := rep.conn.Exec("UPDATE Balance SET balance = balance + $1::integer WHERE userID = $2::integer;", balance, userID)

	if errUpdate != nil {
		return fmt.Errorf("ошибка записи в базу")
	}

	return nil
}

// UpdateUserBalance - Добавление средств на сче
func (rep *BalanceRepository) UpdateUserBalance(userID, balance int) error {

	tx, _ := rep.conn.Begin()

	errTrans := rep.addTransaction(userID, balance, 0)
	errUpdate := rep.updateUserBalance(userID, balance)

	if errTrans != nil || errUpdate != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка записи в базу")
	}

	tx.Commit()
	return nil
}

// Снятие денег со счета
func (rep *BalanceRepository) cashWithdrawal(userID, amount int) error {

	check := rep.balanceCheck(userID, amount)
	if check != true {
		return fmt.Errorf("не достаточно средств")
	}
	_, errBalance := rep.conn.Exec("UPDATE Balance SET balance = balance - $1::integer WHERE userID = $2::integer;", amount, userID)

	if errBalance != nil {
		log.Printf("Данные в базе не обновились: %e", errBalance)
		return errBalance
	}

	return nil
}

// Добавление средств на счет для резервирования
func (rep *BalanceRepository) reserve(userID, amount int) error {

	_, errReserve := rep.conn.Exec("INSERT INTO Reserved (userID, balance) VALUES ($1, $2);", userID, amount)
	if errReserve != nil {
		log.Printf("Данные в базу не записались: %e", errReserve)
		return errReserve
	}
	return nil
}

// Reserve - Резервирование средств
func (rep *BalanceRepository) Reserve(userID, amount int) error {

	if !rep.balanceCheck(userID, amount) {
		return fmt.Errorf("не достаточно средств")
	}

	tx, _ := rep.conn.Begin()

	errWithdrawal := rep.cashWithdrawal(userID, amount)
	errReserve := rep.reserve(userID, amount)
	if errWithdrawal != nil || errReserve != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка записи в базу")
	}

	tx.Commit()
	return nil
}

// GetBalance - Получение баланса. Не нужна транзакция
func (rep *BalanceRepository) GetBalance(userID int) int {

	row := rep.conn.QueryRow(
		"Select sum(balance) FROM Balance Where userid = $1::integer", userID)

	var balance int
	row.Scan(&balance)

	return balance
}

// Узнаем сумму резерва
func (rep *BalanceRepository) getReserve(userID int) int {

	row := rep.conn.QueryRow(
		"Select sum(balance) FROM Reserved Where userid = $1::integer", userID)

	var reserve int
	row.Scan(&reserve)

	return reserve
}

func (rep *BalanceRepository) deleteReserve(userID int) error {

	_, err := rep.conn.Exec("DELETE FROM Reserved Where userid = $1::integer", userID)
	if err != nil {
		log.Printf("Данные в базу не записались: %e", err)
		return err
	}
	return nil
}

// Debit - Списание денег с резервного счета
func (rep *BalanceRepository) Debit(userID int) error {

	amount := rep.getReserve(userID)

	tx, _ := rep.conn.Begin()

	errTrans := rep.addTransaction(userID, -amount, 1)
	errDelete := rep.deleteReserve(userID)

	if errTrans != nil || errDelete != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка записи в базу")
	}

	tx.Commit()
	return nil
}

// Refund - Возврат средств на основной счет
func (rep *BalanceRepository) Refund(userID int) error {

	amount := rep.getReserve(userID)

	tx, _ := rep.conn.Begin()

	errUpdate := rep.updateUserBalance(userID, amount)
	errDelete := rep.deleteReserve(userID)

	if errUpdate != nil || errDelete != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка записи в базу")
	}

	tx.Commit()
	return nil
}

// TransferMoney - Перевод от user к user
func (rep *BalanceRepository) TransferMoney(userID1, userID2, sum int) error {

	check := rep.balanceCheck(userID1, sum)
	if check != true {
		return fmt.Errorf("не достаточно средств")
	}

	tx, _ := rep.conn.Begin()

	errTrans1 := rep.addTransaction(userID1, -sum, 1)
	errWithdrawal := rep.cashWithdrawal(userID1, sum)

	errUpdate := rep.UpdateUserBalance(userID2, sum)

	if errTrans1 != nil || errWithdrawal != nil || errUpdate != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка записи в базу")
	}

	tx.Commit()
	return nil

}

func (rep *BalanceRepository) GetTransactionsReport(userID int) error {

	file, _ := os.Create("transactions_report.csv")

	defer file.Close()

	file.WriteString(fmt.Sprint("id_transaction", ";",
		"type_transaction", ";",
		"userID", ";",
		"amount", ";",
		"time_oper", ";\n"))

	rows, err := rep.conn.Query("SELECT * FROM Transactions WHERE userID  = $1", userID)
	if err != nil {
		return err
	}

	var row struct {
		id_transaction   int
		type_transaction int
		userID           int
		amount           int
		time_oper        time.Time
	}

	for rows.Next() {
		err := rows.Scan(&row.id_transaction,
			&row.type_transaction,
			&row.userID,
			&row.amount,
			&row.time_oper)

		if err != nil {
			return err
		}

		file.WriteString(fmt.Sprint(row.id_transaction, ";",
			row.type_transaction, ";",
			row.userID, ";",
			row.amount, ";",
			row.time_oper, ";\n"))
	}

	return nil
}

func (rep *BalanceRepository) GetAccountingReport(month, year int) error {

	file, _ := os.Create("accounting_report.csv")

	defer file.Close()

	file.WriteString(fmt.Sprint("userID", ";",
		"sum", ";\n"))

	var row struct {
		userID int
		sum    int
	}

	rows, err := rep.conn.Query("SELECT userid, sum (amount) FROM public.transactions Where Extract (Month from time_oper) = $1 and Extract (Year from time_oper) = $2 Group By userid;", month, year)
	if err != nil {
		return err
	}

	for rows.Next() {
		err := rows.Scan(&row.userID,
			&row.sum)

		if err != nil {
			return err
		}

		file.WriteString(fmt.Sprint(row.userID, ";",
			row.sum, ";\n"))
	}

	return nil
}
