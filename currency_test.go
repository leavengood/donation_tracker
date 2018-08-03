package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	var expected float32 = 1.2345
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"base":"EUR","date":"2014-09-27","rates":{"USD":%v}}`, expected)
	}))
	defer ts.Close()

	exchangeRateUrl = ts.URL

	rate, err := GetExchangeRate("?fake-access-key")
	if err != nil {
		t.Errorf("Received unexpected error: %v", err)
	}
	if rate != expected {
		t.Errorf("Expected: %v, but got: %v", expected, rate)
	}
}
