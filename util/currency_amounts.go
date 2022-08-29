package util

import "fmt"

// Use a map because of multiple currencies
type CurrencyAmounts map[string]float32

func (ca1 CurrencyAmounts) Add(ca2 CurrencyAmounts) CurrencyAmounts {
	result := make(CurrencyAmounts)

	// Be sure to get all the currencies from each map
	for _, ca := range []CurrencyAmounts{ca1, ca2} {
		for currency := range ca {
			result[currency] = ca1[currency] + ca2[currency]
		}
	}

	return result
}

// func (c CurrencyAmounts) Sub(c2 CurrencyAmounts) CurrencyAmounts {
// 	result := make(CurrencyAmounts)

// 	for currency, amt := range c {
// 		result[currency] = amt - c2[currency]
// 	}

// 	return result
// }

func (ca CurrencyAmounts) GrandTotal(eurToUsdRate float32) float32 {
	return ca["USD"] + ca["EUR"]*eurToUsdRate
}

func (ca CurrencyAmounts) String() string {
	if ca["USD"] > 0 && ca["EUR"] == 0 {
		return fmt.Sprintf("[USD: %0.02f]", ca["USD"])
	}

	if ca["USD"] == 0 && ca["EUR"] > 0 {
		return fmt.Sprintf("[EUR: %0.02f]", ca["EUR"])
	}

	return fmt.Sprintf("[USD: %0.02f, EUR: %0.02f]", ca["USD"], ca["EUR"])
}
