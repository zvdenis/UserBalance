package Services

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

//Отвечает за логику баланса пользователей
type BalanceManager struct {
	DB *sql.DB
}

//Добавляет сумму пользователю
func (balanceManager BalanceManager) AddMoney(userID int, value int, message string) error {
	_, err := balanceManager.DB.Exec("INSERT INTO balance.balance (id, money) VALUES(?, ?) ON DUPLICATE KEY UPDATE money= money + ?", userID, value, value)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = balanceManager.DB.Exec("INSERT INTO balance.transactions (user_id, message, money) VALUES(?, ?, ?)", userID, fmt.Sprintf("+%d : %s", value, message), value)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

//Возвращает баланс пользователся
func (balanceManager BalanceManager) GetUserMoney(userID int) (int, error) {
	rows, err := balanceManager.DB.Query("SELECT (money) FROM balance.balance WHERE id = ?", userID)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	var money int

	if !rows.Next() {
		return -1, errors.New("no such ID")
	}
	rows.Scan(&money)
	return money, nil
}

//Возвращает баланс пользователся в валюте с заданным кодом
func (balanceManager BalanceManager) GetConvertedUserMoney(userID int, code string) (string, error) {
	money, err := balanceManager.GetUserMoney(userID)
	if err != nil {
		return "", err
	}
	coefficient, err := GetCurrency(code)
	if err != nil {
		return "", err
	}
	convertedMoney := float64(money) * coefficient
	return fmt.Sprintf("%.2f", convertedMoney), nil
}

//Снимает с баланса пользователя заданную сумму
func (balanceManager BalanceManager) DebitMoney(userID int, value int, message string) error {
	money, err := balanceManager.GetUserMoney(userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if value > money {
		return errors.New("not enough money")
	}
	_, err = balanceManager.DB.Query("UPDATE balance.balance SET money = money - ? WHERE id = ?", value, userID)
	if err != nil {
		return err
	}

	_, err = balanceManager.DB.Exec("INSERT INTO balance.transactions (user_id, message, money) VALUES(?, ?, ?)", userID, fmt.Sprintf("-%d : %s", value, message), -money)
	if err != nil {
		fmt.Println(err)
	}


	return nil
}

//переводит деньги с одного счета на другой
func (balanceManager BalanceManager) TransferMoney(userID int, secondUserID int, value int, message string) error {
	money, err := balanceManager.GetUserMoney(userID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if value > money {
		return errors.New("not enough money")
	}
	err = balanceManager.DebitMoney(userID, value, "transfer to id" + strconv.Itoa(secondUserID) + "   Message: " + message)
	if err != nil {
		return err
	}
	err = balanceManager.AddMoney(secondUserID, value, "transfer from id" + strconv.Itoa(userID) + "   Message: " + message)
	if err != nil {
		return err
	}

	return nil
}

//Возвращает список тразнакций пользователя
func (balanceManager BalanceManager) GetUserHistory(userID int) ([]string, error) {
	list := []string{}
	rows, err := balanceManager.DB.Query("SELECT (message) FROM balance.transactions WHERE user_id = ?" , userID)
	if err != nil {
		fmt.Println(err)
		return list, err
	}

	var row string
	for rows.Next(){
		rows.Scan(&row)
		list = append(list, row)
	}
	return list, nil
}
