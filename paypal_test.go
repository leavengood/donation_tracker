package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//==============================================================================
// PayPalTxns.TotalByCurrency
//==============================================================================

func CallTotalByCurrencyWith(txns ...*PayPalTxn) CurrencyCount {
	testing := PayPalTxns{}

	testing = append(testing, txns...)

	return testing.TotalByCurrency()
}

func TestTotalByCurrencyWithEmptyArray(t *testing.T) {
	total := CallTotalByCurrencyWith()

	assert.Equal(t, 0, total["USD"])
}

func TestTotalByCurrencyWithOneItem(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
	)

	assert.Equal(t, 5.43, total["USD"])
}

func TestTotalByCurrencyWithTwoItemsInSameCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 1.25, CurrencyCode: "USD"},
	)

	assert.Equal(t, 6.68, total["USD"])
}

func TestTotalByCurrencyWithTwoItemsInDifferentCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 1.25, CurrencyCode: "EUR"},
	)

	assert.Equal(t, 5.43, total["USD"])
	assert.Equal(t, 1.25, total["EUR"])
}

func TestTotalByCurrencyWithMultipleItemsInDifferentCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 2.12, CurrencyCode: "EUR"},
		&PayPalTxn{Amt: 7.89, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 1.25, CurrencyCode: "EUR"},
	)

	assert.Equal(t, 13.32, total["USD"])
	assert.Equal(t, 3.37, total["EUR"])
}

//==============================================================================
