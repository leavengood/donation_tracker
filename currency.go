package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const FixerIoUrl = "http://api.fixer.io/latest?symbols=USD"

type FixerIoResponse struct {
	Base  string
	Date  string
	Rates map[string]float32
}

var exchangeRateUrl = FixerIoUrl

func GetExchangeRate() (float32, error) {
	resp, err := http.Get(exchangeRateUrl)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fir := new(FixerIoResponse)
	err = json.Unmarshal(body, fir)
	if err != nil {
		return 0, err
	}
	return fir.Rates["USD"], nil
}
