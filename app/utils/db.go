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
		output += fmt.Sprintf("| %s      |%s %.2f       |\n", value, rate.GetCurrencySymbols(value), result)
	}

	return output
}

func (ds *DatabaseService) GetTransactionId(messageId int) (int, error) {
	query := "select transaction_id from payload WHERE (data->>'message_id')::int = $1;"
	rows, err := ds.db.Query(context.Background(), query, messageId)
	if err != nil {
		fmt.Println(err)
	}

	var trId int
	for rows.Next() {
		if err := rows.Scan(&trId); err != nil {
			return 0, fmt.Errorf("Ошибка при чтении строки: %v\n", err)
		}
	}

	return trId, nil
}

func (ds *DatabaseService) UpdateTransaction(info TransactionType, transactionID int) error {
	fmt.Println(info, transactionID)
	query := "update money set amount = $1, currency = $2, type = $3, name = $4, description = $5, usd_rate = $6, updated_at = now() where id = $7;"
	_, err := ds.db.Exec(context.Background(), query, info.Amount, info.Currency, info.Type, info.Name, info.Description, 1, transactionID)
	return err
}

func (ds *DatabaseService) SavePayload(transaction NewPayloadMessage, transactionID int) {
	query := "insert into payload (transaction_id, data) values ($1, $2);"
	ds.db.Exec(context.Background(), query, transactionID, transaction)
}
func (ds *DatabaseService) UpdatePayload(transaction NewPayloadMessage, transactionID int) {
	query := "UPDATE payload SET data = jsonb_set(data, '{update}', $2) WHERE transaction_id = $1;"
	ds.db.Exec(context.Background(), query, transactionID, transaction)
}

func (ds *DatabaseService) SaveInfo(text string, rate *RateService) (int, error) {
	rateMap := rate.rateMap
	parse, parseErr := ParseText(text)
	if parseErr != nil {
		return 0, parseErr
	}

	info := rateMap[strings.ToLower(parse.Currency)]
	usdRate := 1 / info

	var id int
	query := "insert into money (amount, currency, type, name, description, usd_rate) values ($1,$2, $3,$4, $5, $6) returning id;"
	err := ds.db.QueryRow(context.Background(), query, parse.Amount, parse.Currency, parse.Type, parse.Name, parse.Description, usdRate).Scan(&id)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, fmt.Errorf("Transaction not save, please try again %s", "")
	}

	return id, nil
}
