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

func (ds *DatabaseService) GetAll(rate *RateService) string {

	rows, err := ds.db.Query(context.Background(), "select amount, currency, usd_rate from money")
	if err != nil {
		fmt.Println(err)
	}

	var currencyAccountMap = make(map[string]string)
	var bank float64

	for rows.Next() {

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

	finalBank := fmt.Sprintf("%.2f", bank)
	output += "# Your Balances =" + "  $" + finalBank + "\n"
	output += "| Currency 	| Balance |\n"
	output += "|------------|---------|\n"

	//// Форматированный вывод остатка
	for _, value := range currencyAccountMap {
		result, _ := rate.ConvertCurrency(bank, "usd", strings.ToLower(value))
		output += fmt.Sprintf("| %s      |%s %.2f       |\n", value, rate.GetCurrencySymbols(value), result)
	}

	return output
}

func (ds *DatabaseService) SaveInfo(text string, rate *RateService) (string, error) {
	rateMap := rate.rateMap
	parse, parseErr := ParseText(text)
	if parseErr != nil {
		return "", parseErr
	}

	info := rateMap[strings.ToLower(parse.Currency)]
	usdRate := 1 / info

	_, err := ds.db.Exec(context.Background(), "insert into money (amount, currency, type, name, description, usd_rate) values ($1,$2, $3,$4, $5, $6) ", parse.Amount, parse.Currency, parse.Type, parse.Name, parse.Description, usdRate)
	if err != nil {
		return "", err
	}

	return ds.GetAll(rate), nil
}
