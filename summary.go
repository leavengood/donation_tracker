package main

// Use a map because of multiple currencies
type CurrencyCount map[string]float32

func (c CurrencyCount) GrandTotal(eurToUsdRate float32) float32 {
	return c["USD"] + c["EUR"]*eurToUsdRate
}
