package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//==============================================================================
// PayPalTxns.TotalByCurrency
//==============================================================================

func CallTotalByCurrencyWith(txns ...*PayPalTxn) CurrencyAmounts {
	testing := PayPalTxns{}

	testing = append(testing, txns...)

	return testing.TotalByCurrency()
}

func TestTotalByCurrencyWithEmptyArray(t *testing.T) {
	total := CallTotalByCurrencyWith()

	assert.Equal(t, float32(0), total["USD"])
}

func TestTotalByCurrencyWithOneItem(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
	)

	assert.Equal(t, float32(5.43), total["USD"])
}

func TestTotalByCurrencyWithTwoItemsInSameCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 1.25, CurrencyCode: "USD"},
	)

	assert.Equal(t, float32(6.68), total["USD"])
}

func TestTotalByCurrencyWithTwoItemsInDifferentCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 1.25, CurrencyCode: "EUR"},
	)

	assert.Equal(t, float32(5.43), total["USD"])
	assert.Equal(t, float32(1.25), total["EUR"])
}

func TestTotalByCurrencyWithMultipleItemsInDifferentCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&PayPalTxn{Amt: 5.43, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 2.12, CurrencyCode: "EUR"},
		&PayPalTxn{Amt: 7.89, CurrencyCode: "USD"},
		&PayPalTxn{Amt: 1.25, CurrencyCode: "EUR"},
	)

	assert.Equal(t, float32(13.32), total["USD"])
	assert.Equal(t, float32(3.37), total["EUR"])
}

//==============================================================================
// PayPalTxns.FilterDonations
//==============================================================================

func CallFilterDonationsWith(txns ...*PayPalTxn) PayPalTxns {
	testing := PayPalTxns{}

	testing = append(testing, txns...)

	donations, _ := testing.FilterDonations()

	return donations
}

func TestFilterDonationsWithEmptyArray(t *testing.T) {
	assert.Equal(t, 0, len(CallFilterDonationsWith()))
}

func TestFilterDonationsWithOneItem(t *testing.T) {
	result := CallFilterDonationsWith(
		&PayPalTxn{Amt: 5.43, Type: "Donation"},
	)

	assert.Equal(t, 1, len(result))
}

func TestFilterDonationsWithAPayment(t *testing.T) {
	result := CallFilterDonationsWith(
		&PayPalTxn{Amt: -5.43, Type: "Payment"},
	)

	assert.Equal(t, 0, len(result))
}

func TestFilterDonationsWithOneDonationAndOnePayment(t *testing.T) {
	result := CallFilterDonationsWith(
		&PayPalTxn{Amt: 5.43, Type: "Donation"},
		&PayPalTxn{Amt: -3.12, Type: "Payment"},
	)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, float32(5.43), result[0].Amt)
}

func TestFilterDonationsWithADonationSubscriptionAndPayment(t *testing.T) {
	result := CallFilterDonationsWith(
		&PayPalTxn{Amt: 5.43, Type: "Donation"},
		&PayPalTxn{Amt: -3.12, Type: "Payment"},
		&PayPalTxn{Amt: 2.45, Type: "Payment"},
	)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, float32(5.43), result[0].Amt)
	assert.Equal(t, float32(2.45), result[1].Amt)
}
