package transaction

import (
	"fmt"
	"main/utils"
	"strconv"
	"strings"
)

func isInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func isString(v interface{}) bool {
	_, ok := v.(string)
	return ok
}

type TransactionInfo struct {
	Amount      int     `json:"amount"`
	Currency    string  `json:"currency"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Usd_Rate    float64 `json:"usd_rate"`
	Description *string `json:"description,omitempty"`
}

func getTypeTransaction(amount int) string {
	if amount < 0 {
		return "Withdraw"
	}
	return "Deposit"
}

func ParseText(input string, rate *utils.RateService) (TransactionInfo, error) {
	var result TransactionInfo
	parts := strings.Split(input, " ")
	if len(parts) < 3 {
		return TransactionInfo{}, fmt.Errorf("not enough arguments. For example: <Amount> <Currency> <Type>.<Description>")
	}
	if !isInteger(parts[0]) {
		return TransactionInfo{}, fmt.Errorf("not correct input. For example: <Amount> is number")
	}
	amount, _ := strconv.Atoi(parts[0])
	result.Amount = amount
	result.Type = getTypeTransaction(amount)

	if !isString(parts[1]) {
		return TransactionInfo{}, fmt.Errorf("not correct input. For example: <Currency> is string")
	}
	result.Currency = strings.ToUpper(parts[1])
	info := rate.Get()[strings.ToLower(strings.ToLower(result.Currency))]
	if info == 0 {
		return TransactionInfo{}, fmt.Errorf("not correct input. For example: <Currency> is not correct")
	}
	result.Usd_Rate = 1 / info
	nameTransaction := strings.Join(parts[2:], " ")

	aboutTransaction := strings.Split(nameTransaction, ".")
	result.Name = aboutTransaction[0]
	if len(aboutTransaction) > 2 {
		result.Description = &aboutTransaction[1]
	}
	return result, nil
}
