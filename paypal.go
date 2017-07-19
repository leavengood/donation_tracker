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
	Timestamp     time.Time
	Type          string
	Email         string
	Name          string
	TransactionID string `json:"transaction_id"`
	Status        string
	Amt           float32
	FeeAmt        float32 `json:"fee_amt"`
	NetAmt        float32 `json:"net_amt"`
	CurrencyCode  string  `json:"currency_code"`
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
	return p.Amt > 0 && p.Type == "Payment"
}

func (p *PayPalTxn) IsDonation() bool {
	return p.Amt > 0 && p.Type == "Donation"
}

func (p *PayPalTxn) String() string {
	return fmt.Sprintf("%s %s <%s> %s, %s %0.02f (%0.02f fee) = %0.02f, %s", p.Timestamp, p.Name, p.Email,
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
