package utils

import (
	"fmt"
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

type TransactionType struct {
	Amount      int     `json:"amount"`
	Currency    string  `json:"currency"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description *string `json:"description,omitempty"`
}

func getTypeTransaction(amount int) string {
	if amount < 0 {
		return "Withdraw"
	}
	return "Deposit"
}

func ParseText(input string) (TransactionType, error) {
	var result TransactionType
	parts := strings.Split(input, " ")
	if len(parts) < 3 {
		return TransactionType{}, fmt.Errorf("not enough arguments. For example: <Amount> <Currency> <Type>.<Description>")
	}
	if !isInteger(parts[0]) {
		return TransactionType{}, fmt.Errorf("not correct input. For example: <Amount> is number")
	}
	amount, _ := strconv.Atoi(parts[0])
	result.Amount = amount
	result.Type = getTypeTransaction(amount)

	if !isString(parts[1]) {
		return TransactionType{}, fmt.Errorf("not correct input. For example: <Currency> is string")
	}
	result.Currency = strings.ToUpper(parts[1])
	nameTransaction := strings.Join(parts[2:], " ")

	aboutTransaction := strings.Split(nameTransaction, ".")
	result.Name = aboutTransaction[0]
	if len(aboutTransaction) > 2 {
		result.Description = &aboutTransaction[1]
	}
	return result, nil
}
