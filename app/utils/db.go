package utils

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"strings"
)

type DatabaseService struct {
	db *pgx.Conn
}

func NewDatabaseService() *DatabaseService {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", DbUser, DbPassword, DbHost, DbName)
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v\n", err)
	}

	return &DatabaseService{db: db}
}

func (ds *DatabaseService) Close() {
	ds.db.Close(context.Background())
}

type Info struct {
	Amount   float64
	Currency string
	Rate     float64
}

func (ds *DatabaseService) GetMoney(rate *RateService) string {

	rows, err := ds.db.Query(context.Background(), "select amount, currency, usd_rate from money")
	if err != nil {
		fmt.Println(err)
	}

	var currencyAccountMap = make(map[string]string)
	var bank float64

	for rows.Next() {
		fmt.Println(rows.Values())
		info := Info{}
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
	fmt.Println(currencyAccountMap)
	finalBank := fmt.Sprintf("%.2f", bank)
	output += "# Your Balances =" + "  $" + finalBank + "\n"
	output += "| Currency 	| Balance |\n"
	output += "|------------|---------|\n"

	//// Форматированный вывод остатка
	for _, value := range currencyAccountMap {
		result, _ := rate.ConvertCurrency(bank, "usd", strings.ToLower(value))
		output += fmt.Sprintf("| %s      |$%.2f       |\n", value, result)
	}

	return output
}

func (ds *DatabaseService) SaveInfo(text string, rate *RateService) string {

	rateMap := rate.rateMap

	var typeName = "Deposit"
	parts := strings.Split(text, " ")
	var description = ""
	amount, currency, name := parts[0], parts[1], parts[2]

	if len(parts) > 3 {
		description = parts[3]
	}

	if amount[0] == '-' {
		typeName = "Withdraw"
	}

	info := rateMap[strings.ToLower(currency)]
	usdRate := 1 / info

	_, err := ds.db.Exec(context.Background(), "insert into money (amount, currency, type, name, description, usd_rate) values ($1,$2, $3,$4, $5, $6) ", amount, currency, typeName, name, description, usdRate)
	if err != nil {
		return err.Error()
	}

	return ds.GetMoney(rate)
}
