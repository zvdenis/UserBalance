package Services

import (
	"encoding/json"
	"net/http"
)

const url = "https://api.exchangeratesapi.io/latest?base=RUB"

type currencyList struct{
	Base string `json:"base"`
	Date string `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func GetCurrency(code string) (float64, error) {
	resp, err := http.Get(url)
	if err != nil{
		return 0, err
	}
	var currency currencyList
	err = json.NewDecoder(resp.Body).Decode(&currency)
	if err != nil{
		return 0, err
	}
	return currency.Rates[code], nil
}
