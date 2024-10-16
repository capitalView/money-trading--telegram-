package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetRequest(date string, currency string) (map[string]float64, error) {
	if date == "" {
		date = "latest"
	}
	if currency == "" {
		currency = "usd"
	}
	url := fmt.Sprintf("https://%s.currency-api.pages.dev/v1/currencies/%s.json", date, currency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	d, _ := ParseJson(string(body))

	var mapData = make(map[string]float64)

	dataMap := d.(map[string]interface{})
	for key, value := range dataMap {
		if key == "date" {
			continue
		}
		currencyData := value.(map[string]interface{})
		for key, value := range currencyData {
			mapData[key] = value.(float64)
		}
	}

	return mapData, nil
}

type RateService struct {
	rateMap map[string]float64
}

func NewRateService() (*RateService, error) {
	rateMap, err := GetRequest("", "")
	if err != nil {
		return nil, err
	}
	return &RateService{rateMap: rateMap}, nil
}

func (rs *RateService) UpdateRates() {
	rateMap, _ := GetRequest("", "")
	rs.rateMap = rateMap
}

func (rs *RateService) ConvertCurrency(amount float64, fromCurrency string, toCurrency string) (float64, error) {
	fromRate, fromExists := rs.rateMap[fromCurrency]
	toRate, toExists := rs.rateMap[toCurrency]

	// Проверяем, существуют ли курсы для обеих валют
	if !fromExists || !toExists {
		return 0, fmt.Errorf("currency not found: %s or %s", fromCurrency, toCurrency)
	}

	// Конвертация суммы: сначала в USD, затем в целевую валюту
	amountInUSD := amount / fromRate
	convertedAmount := amountInUSD * toRate

	return convertedAmount, nil
}
