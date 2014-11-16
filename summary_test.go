package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//==============================================================================
// CurrencyCount.GrandTotal
//==============================================================================

const eurToUsdRate = 1.25

func TestGrandTotalWithEmptyMap(t *testing.T) {
	cc := make(CurrencyCount)

	assert.Equal(t, 0, cc.GrandTotal(eurToUsdRate))
}

func TestGrandTotalWithJustUSD(t *testing.T) {
	cc := CurrencyCount{
		"USD": 34.56,
	}

	assert.Equal(t, 34.56, cc.GrandTotal(eurToUsdRate))
}

func TestGrandTotalWithJustUSDAndEUR(t *testing.T) {
	cc := CurrencyCount{
		"USD": 34.56,
		"EUR": 10.00,
	}

	assert.Equal(t, 47.06, cc.GrandTotal(eurToUsdRate))
}

//==============================================================================
