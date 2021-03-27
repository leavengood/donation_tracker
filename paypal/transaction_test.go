package paypal

import (
	"testing"

	"github.com/leavengood/donation_tracker/util"
	"github.com/stretchr/testify/assert"
)

//==============================================================================
// TotalByCurrency
//==============================================================================

func CallTotalByCurrencyWith(txns ...*Transaction) util.CurrencyAmounts {
	testing := Transactions{}

	testing = append(testing, txns...)

	return testing.TotalByCurrency()
}

func TestTotalByCurrencyWithEmptyArray(t *testing.T) {
	total := CallTotalByCurrencyWith()

	assert.Equal(t, float32(0), total["USD"])
}

func TestTotalByCurrencyWithOneItem(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&Transaction{Amt: 5.43, CurrencyCode: "USD"},
	)

	assert.Equal(t, float32(5.43), total["USD"])
}

func TestTotalByCurrencyWithTwoItemsInSameCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&Transaction{Amt: 5.43, CurrencyCode: "USD"},
		&Transaction{Amt: 1.25, CurrencyCode: "USD"},
	)

	assert.Equal(t, float32(6.68), total["USD"])
}

func TestTotalByCurrencyWithTwoItemsInDifferentCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&Transaction{Amt: 5.43, CurrencyCode: "USD"},
		&Transaction{Amt: 1.25, CurrencyCode: "EUR"},
	)

	assert.Equal(t, float32(5.43), total["USD"])
	assert.Equal(t, float32(1.25), total["EUR"])
}

func TestTotalByCurrencyWithMultipleItemsInDifferentCurrency(t *testing.T) {
	total := CallTotalByCurrencyWith(
		&Transaction{Amt: 5.43, CurrencyCode: "USD"},
		&Transaction{Amt: 2.12, CurrencyCode: "EUR"},
		&Transaction{Amt: 7.89, CurrencyCode: "USD"},
		&Transaction{Amt: 1.25, CurrencyCode: "EUR"},
	)

	assert.Equal(t, float32(13.32), total["USD"])
	assert.Equal(t, float32(3.37), total["EUR"])
}

//==============================================================================
// FilterDonations
//==============================================================================

func CallFilterDonationsWith(txns ...*Transaction) Transactions {
	testing := Transactions{}

	testing = append(testing, txns...)

	donations, _ := testing.FilterDonations()

	return donations
}

func TestFilterDonationsWithEmptyArray(t *testing.T) {
	assert.Equal(t, 0, len(CallFilterDonationsWith()))
}

func TestFilterDonationsWithOneItem(t *testing.T) {
	result := CallFilterDonationsWith(
		&Transaction{Amt: 5.43, Type: "Donation"},
	)

	assert.Equal(t, 1, len(result))
}

func TestFilterDonationsWithAPayment(t *testing.T) {
	result := CallFilterDonationsWith(
		&Transaction{Amt: -5.43, Type: "Payment"},
	)

	assert.Equal(t, 0, len(result))
}

func TestFilterDonationsWithOneDonationAndOnePayment(t *testing.T) {
	result := CallFilterDonationsWith(
		&Transaction{Amt: 5.43, Type: "Donation"},
		&Transaction{Amt: -3.12, Type: "Payment"},
	)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, float32(5.43), result[0].Amt)
}

func TestFilterDonationsWithADonationSubscriptionAndPayment(t *testing.T) {
	result := CallFilterDonationsWith(
		&Transaction{Amt: 5.43, Type: "Donation"},
		&Transaction{Amt: -3.12, Type: "Payment"},
		&Transaction{Amt: 2.45, Type: "Payment"},
	)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, float32(5.43), result[0].Amt)
	assert.Equal(t, float32(2.45), result[1].Amt)
}
