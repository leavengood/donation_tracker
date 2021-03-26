package main

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const PayPalDateFormat = "2006-01-02T15:04:05Z"

type PayPalTxn struct {
	Timestamp     time.Time `json:"timestamp,omitempty"`
	Type          string    `json:"type,omitempty"`
	Email         string    `json:"email,omitempty"`
	Name          string    `json:"name,omitempty"`
	TransactionID string    `json:"transaction_id,omitempty"`
	Status        string    `json:"status,omitempty"`
	Amt           float32   `json:"amt,omitempty"`
	FeeAmt        float32   `json:"fee_amt,omitempty"`
	NetAmt        float32   `json:"net_amt,omitempty"`
	CurrencyCode  string    `json:"currency_code,omitempty"`
}

func NewPayPalTxn(tran map[string]string) *PayPalTxn {
	result := new(PayPalTxn)

	elem := reflect.ValueOf(result).Elem()

	for name, value := range tran {
		field := elem.FieldByNameFunc(func(n string) bool {
			return name == strings.ToUpper(n)
		})
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Float32:
				f, _ := strconv.ParseFloat(value, 32)
				field.SetFloat(f)
			default:
				// its the Timestamp field
				timestamp, _ := time.Parse(PayPalDateFormat, value)
				field.Set(reflect.ValueOf(timestamp))
			}
		}
	}

	return result
}

func (p *PayPalTxn) IsSubscription() bool {
	return p.Amt > 0 &&
		(p.Type == "Payment" || p.Type == "Recurring Payment")
}

func (p *PayPalTxn) IsDonation() bool {
	return p.Amt > 0 && p.Type == "Donation"
}

func (p *PayPalTxn) String() string {
	tsStr := FormatDateTime(p.Timestamp)

	// For subscription changes to display nicely
	if (p.Type == "Recurring Payment" && p.Amt == 0 && p.FeeAmt == 0) ||
		p.Type == "Subscription Cancellation" {
		color := red

		if p.Status == "Created" {
			color = green
		}

		return Colorize(color, fmt.Sprintf("%s: %s %s a subscription",
			tsStr, p.Name, strings.ToLower(p.Status)))
	}

	return fmt.Sprintf("%s: %s <%s> %s, %s %0.02f (%0.02f fee) = %0.02f, %s", tsStr, p.Name, p.Email,
		p.Type, p.CurrencyCode, p.Amt, p.FeeAmt, p.NetAmt, p.Status)
}

type PayPalTxns []*PayPalTxn

func PayPalTxnsFromNvp(nvp *NvpResult) PayPalTxns {
	result := make(PayPalTxns, len(nvp.List))

	for i, item := range nvp.List {
		result[i] = NewPayPalTxn(item)
		// fmt.Println(result[i])
	}

	return result
}

// Sorting

func (p PayPalTxns) Len() int      { return len(p) }
func (p PayPalTxns) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type ByDate struct{ PayPalTxns }

func (s ByDate) Less(i, j int) bool {
	return s.PayPalTxns[i].Timestamp.Before(s.PayPalTxns[j].Timestamp)
}

func (p PayPalTxns) Sort() {
	sort.Sort(ByDate{p})
}

func (p PayPalTxns) TotalByCurrency() CurrencyAmounts {
	result := make(CurrencyAmounts)

	for _, txn := range p {
		result[txn.CurrencyCode] += txn.Amt
	}

	return result
}

func (p PayPalTxns) FilterDonations() (PayPalTxns, PayPalTxns) {
	donations := make(PayPalTxns, 0, len(p))
	other := make(PayPalTxns, 0, len(p))

	for _, item := range p {
		if item.IsDonation() || item.IsSubscription() {
			donations = append(donations, item)
		} else {
			other = append(other, item)
		}
	}

	return donations, other
}

func (p PayPalTxns) Summarize() MonthlySummaries {
	result := make(MonthlySummaries)

	for _, item := range p {
		_, month, _ := item.Timestamp.Date()
		summary := result.ForMonth(month)

		if item.IsDonation() {
			summary.AddOneTime(item.Amt, item.FeeAmt, item.CurrencyCode)
		} else if item.IsSubscription() {
			summary.AddSubscription(item.Amt, item.FeeAmt, item.CurrencyCode)
		}
	}

	return result
}

func (p PayPalTxns) Merge(other PayPalTxns) PayPalTxns {
	result := make(PayPalTxns, 0, len(p)+len(other))
	tranIDs := map[string]bool{}

	// Add everything in our own list, tracking transaction IDs
	for _, item := range p {
		tranIDs[item.TransactionID] = true
		result = append(result, item)
	}

	// Add anything new not already in the list
	for _, item := range other {
		if _, found := tranIDs[item.TransactionID]; !found {
			result = append(result, item)
		}
	}

	return result
}
