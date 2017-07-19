package main

import (
	"encoding/csv"
	"os"
	"reflect"
	"strconv"
	"time"
)

const CsvDateFormat = "2006-01-02"

type Transaction struct {
	Date         time.Time
	Type         string
	Name         string
	Email        string
	Amt          float32
	FeeAmt       float32
	CurrencyCode string
}

func NewTransaction(row []string, headers []string) *Transaction {
	result := new(Transaction)

	elem := reflect.ValueOf(result).Elem()

	for i, header := range headers {
		field := elem.FieldByName(header)

		if field.IsValid() && field.CanSet() {
			value := row[i]

			switch field.Kind() {
			case reflect.String:
				field.SetString(value)
			case reflect.Float32:
				f, _ := strconv.ParseFloat(value, 32)
				field.SetFloat(f)
			default:
				// its the Date field
				date, _ := time.Parse(CsvDateFormat, value)
				field.Set(reflect.ValueOf(date))
			}
		}
	}

	return result
}

func (t *Transaction) NetAmt() float32 {
	return t.Amt - t.FeeAmt
}

func ReadCsv(name string) ([][]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return [][]string{}, err
	}

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return rows, nil
}

func TransactionsFromCsv(name string) ([]*Transaction, error) {
	rows, err := ReadCsv(name)
	if err != nil {
		return []*Transaction{}, err
	}

	ts := make([]*Transaction, len(rows)-1)

	headers := rows[0]

	for i := 1; i < len(rows); i++ {
		ts[i-1] = NewTransaction(rows[i], headers)
	}

	return ts, nil
}
