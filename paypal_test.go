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
// PayPalTxns.OnlyDonations
//==============================================================================

func CallOnlyDonationsWith(txns ...*PayPalTxn) PayPalTxns {
	testing := PayPalTxns{}

	testing = append(testing, txns...)

	return testing.OnlyDonations()
}

func TestOnlyDonationsWithEmptyArray(t *testing.T) {
	assert.Equal(t, 0, len(CallOnlyDonationsWith()))
}

func TestOnlyDonationsWithOneItem(t *testing.T) {
	result := CallOnlyDonationsWith(
		&PayPalTxn{Amt: 5.43, Type: "Donation"},
	)

	assert.Equal(t, 1, len(result))
}

func TestOnlyDonationsWithAPayment(t *testing.T) {
	result := CallOnlyDonationsWith(
		&PayPalTxn{Amt: -5.43, Type: "Payment"},
	)

	assert.Equal(t, 0, len(result))
}

func TestOnlyDonationsWithOneDonationAndOnePayment(t *testing.T) {
	result := CallOnlyDonationsWith(
		&PayPalTxn{Amt: 5.43, Type: "Donation"},
		&PayPalTxn{Amt: -3.12, Type: "Payment"},
	)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 5.43, result[0].Amt)
}

func TestOnlyDonationsWithADonationSubscriptionAndPayment(t *testing.T) {
	result := CallOnlyDonationsWith(
		&PayPalTxn{Amt: 5.43, Type: "Donation"},
		&PayPalTxn{Amt: -3.12, Type: "Payment"},
		&PayPalTxn{Amt: 2.45, Type: "Payment"},
	)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, 5.43, result[0].Amt)
	assert.Equal(t, 2.45, result[1].Amt)
}
