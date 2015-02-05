package main

// Use a map because of multiple currencies
type CurrencyAmounts map[string]float32

func (c CurrencyAmounts) GrandTotal(eurToUsdRate float32) float32 {
	return c["USD"] + c["EUR"]*eurToUsdRate
}
