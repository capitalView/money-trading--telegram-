package transaction

import (
	"fmt"
	"log"
	"main/db"
	"main/utils"
	"strings"
)

type Transaction struct {
	*db.Database
	rate *utils.RateService
}

func NewTransaction(db *db.Database, rate *utils.RateService) *Transaction {
	return &Transaction{Database: db, rate: rate}
}

func (data *Transaction) Insert(t TransactionInfo) (int, error) {
	var id int
	err := data.QueryRow(InsertTransaction, t.Amount, t.Currency, t.Type, t.Name, t.Description, t.Usd_Rate).Scan(&id)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, fmt.Errorf("Transaction not save, please try again")
	}

	return id, nil
}

func (data *Transaction) Update(transactionID int, info TransactionInfo) error {
	err := data.Execute(UpdateTransaction, info.Amount, info.Currency, info.Type, info.Name, info.Description, info.Usd_Rate, transactionID)
	return err
}

func (data *Transaction) GetAll() string {
	rate := data.rate
	rows, err := data.Query(SelectTransactions)
	if err != nil {
		fmt.Println(err)
	}

	var currencyAccountMap = make(map[string]string)
	var bank float64

	for rows.Next() {

		info := &Info{}
		err := rows.Scan(&info.Amount, &info.Currency, &info.Rate)
		if err != nil {
			log.Fatalf("Ошибка при чтении строки: %v\n", err)
		}

		_, has := currencyAccountMap[info.Currency]
		if has == false {
			currencyAccountMap[info.Currency] = info.Currency
		}
		bank += info.Amount * info.Rate
	}

	var output string

	finalBank := fmt.Sprintf("%.2f", bank)
	output += "# Your Balances =" + "  $" + finalBank + "\n"
	output += "| Currency 	| Balance |\n"
	output += "|------------|---------|\n"

	//// Форматированный вывод остатка
	for _, value := range currencyAccountMap {
		result, _ := rate.ConvertCurrency(bank, "usd", strings.ToLower(value))
		fmt.Println(value, result)
		output += fmt.Sprintf("| %s      |%s %.2f       |\n", value, rate.GetCurrencySymbols(value), result)
	}

	return output
}
