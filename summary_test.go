package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//==============================================================================
// CurrencyAmounts.Add
//==============================================================================

func TestAddWithAllCurrencies(t *testing.T) {
	ca := CurrencyAmounts{
		"USD": 12.34,
		"EUR": 10.00,
	}
	ca2 := CurrencyAmounts{
		"USD": 23.45,
		"EUR": 20.00,
	}

	result := ca.Add(ca2)

	assert.Equal(t, 35.79, result["USD"])
	assert.Equal(t, 30.00, result["EUR"])
}

func TestAddWithSomeMissingCurrencies(t *testing.T) {
	ca := CurrencyAmounts{
		"EUR": 10.00,
	}
	ca2 := CurrencyAmounts{
		"USD": 23.45,
	}

	result := ca.Add(ca2)

	assert.Equal(t, 23.45, result["USD"])
	assert.Equal(t, 10.00, result["EUR"])
}

//==============================================================================
// CurrencyAmounts.Sub
//==============================================================================

// func TestSub(t *testing.T) {
// 	ca := CurrencyAmounts{
// 		"USD": 12.34,
// 		"EUR": 10.00,
// 	}
// 	ca2 := CurrencyAmounts{
// 		"USD": 23.45,
// 		"EUR": 20.00,
// 	}

// 	result := ca2.Sub(ca)

// 	assert.InDelta(t, 11.11, result["USD"], 0.0001)
// 	assert.Equal(t, 10.00, result["EUR"])
// }

//==============================================================================
// CurrencyAmounts.GrandTotal
//==============================================================================

const eurToUsdRate = 1.25

func TestGrandTotalWithEmptyMap(t *testing.T) {
	ca := make(CurrencyAmounts)

	assert.Equal(t, 0, ca.GrandTotal(eurToUsdRate))
}

func TestGrandTotalWithJustUSD(t *testing.T) {
	ca := CurrencyAmounts{
		"USD": 34.56,
	}

	assert.Equal(t, 34.56, ca.GrandTotal(eurToUsdRate))
}

func TestGrandTotalWithJustUSDAndEUR(t *testing.T) {
	ca := CurrencyAmounts{
		"USD": 34.56,
		"EUR": 10.00,
	}

	assert.Equal(t, 47.06, ca.GrandTotal(eurToUsdRate))
}

//==============================================================================
