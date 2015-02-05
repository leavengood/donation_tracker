package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//==============================================================================
// CurrencyAmounts.GrandTotal
//==============================================================================

const eurToUsdRate = 1.25

func TestGrandTotalWithEmptyMap(t *testing.T) {
	cc := make(CurrencyAmounts)

	assert.Equal(t, 0, cc.GrandTotal(eurToUsdRate))
}

func TestGrandTotalWithJustUSD(t *testing.T) {
	cc := CurrencyAmounts{
		"USD": 34.56,
	}

	assert.Equal(t, 34.56, cc.GrandTotal(eurToUsdRate))
}

func TestGrandTotalWithJustUSDAndEUR(t *testing.T) {
	cc := CurrencyAmounts{
		"USD": 34.56,
		"EUR": 10.00,
	}

	assert.Equal(t, 47.06, cc.GrandTotal(eurToUsdRate))
}

//==============================================================================
