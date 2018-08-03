package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const FixerIoUrl = "http://data.fixer.io/api/latest?format=1&symbols=USD&access_key="

type FixerIoResponse struct {
	Base  string
	Date  string
	Rates map[string]float32
}

var exchangeRateUrl = FixerIoUrl

func GetExchangeRate(accessKey string) (float32, error) {
	resp, err := http.Get(exchangeRateUrl+accessKey)
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
