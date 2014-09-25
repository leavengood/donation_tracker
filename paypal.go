package main

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

type PayPalTxn struct {
	Timestamp time.Time
	Type string
	Email string
	Name string
	TransactionID string
	Status string
	Amt float32
	FeeAmt float32
	NetAmt float32
	CurrencyCode string
}

type PayPalTxns []*PayPalTxn

func NewPayPalTxn(tran map[string]string) *PayPalTxn {
	result := new(PayPalTxn)

	elem := reflect.ValueOf(result).Elem()

	for name, value := range(tran) {
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
				timestamp, _ := time.Parse("2006-01-02T15:04:05Z", value)
				field.Set(reflect.ValueOf(timestamp))
			}
		}
	}

	return result
}

func PayPalTxnsFromNvp(nvp *NvpResult) PayPalTxns {
	result := make(PayPalTxns, len(nvp.List))

	for i, item := range(nvp.List) {
		result[i] = NewPayPalTxn(item)
		// fmt.Println(result[i])
	}

	return result
}
