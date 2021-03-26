package main

import (
	"path/filepath"
	"fmt"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
)

const dataDir = "data"

// PayPalFileManager manages files containing PayPal transactions fetched from
// the PayPal API.
type PayPalFileManager struct {
	Year int
	Months map[int]PayPalTxns
}

// NewPayPalFileManager will load any PayPal files for the given year from the
// data directory, and can be used to save new transactions there.
func NewPayPalFileManager(year int) (*PayPalFileManager, error) {
	months, err := loadPayPalFiles(dataDir, year)
	if err != nil {
		return nil, err
	}

	return &PayPalFileManager{
		Year: year,
		Months: months,
	}, nil
}

func NewEmptyPayPalFileManager(year int) *PayPalFileManager {
	return &PayPalFileManager{
		Year: year,
		Months: map[int]PayPalTxns{},
	}
}

// GetLatestMonth returns the latest month with transactions loaded by this
// file manager. When managing the current year, this helps determine what new
// data needs to be fetched from the PayPal API.
func (p *PayPalFileManager) GetLatestMonth() int {
	max := 0

	for month := range p.Months {
		if month > max {
			max = month
		}
	}

	return max
}

// GetLatestTransaction will get the most recent transaction from the most
// recent month stored in this file manager.
func (p *PayPalFileManager) GetLatestTransaction() *PayPalTxn {
	txns := p.Months[p.GetLatestMonth()]

	// The transactions should be sorted with the latest last
	if len(txns) > 0 {
		return txns[len(txns) - 1]
	}

	return nil
}

// GetExistingMonths will return all months which are currently stored in this
// file manager.
func (p *PayPalFileManager) GetExistingMonths() []int {
	result := []int{}

	for i := 1; i <= 12; i++ {
		if _, found := p.Months[i]; found {
			result = append(result, i)
		}
	}

	return result
}

// GetMissingMonths will return all months which are not currently stored in
// this file manager.
func (p *PayPalFileManager) GetMissingMonths() []int {
	result := []int{}

	for i := 1; i <= 12; i++ {
		if _, found := p.Months[i]; !found {
			result = append(result, i)
		}
	}

	return result
}

// SaveMonth will save the given transactions to a file for that month, and add
// these transactions to the Months stored in this manager.
func (p *PayPalFileManager) SaveMonth(month int, txns PayPalTxns) error {
	filename := payPalTxnsFileName(p.Year, month)
	if err := savePayPalTxnsToFile(filename, txns); err != nil {
		return err
	}

	p.Months[month] = txns

	return nil
}

func loadPayPalFiles(baseDir string, year int) (map[int]PayPalTxns, error) {
	result := map[int]PayPalTxns{}

	re := regexp.MustCompile(fmt.Sprintf(`paypal-%d-([0-9]{2}).json`, year))

	err := filepath.Walk(baseDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if match := re.FindStringSubmatch(path); match != nil {
					// Get the month from the match
					month, err := strconv.Atoi(match[1])
					if err != nil {
						return err
					}
					// fmt.Printf("Found file for year %d: %s. It is for month %d.\n", year, path, month)
					txns, err := loadPayPalTxnsFromFile(path)
					if err != nil {
						return err
					}
					result[month] = txns
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PayPalTxnsJSON is the structure that will be saved to and loaded from files.
type PayPalTxnsJSON struct {
	Transactions PayPalTxns `json:"transactions"`
}

func loadPayPalTxnsFromFile(filename string) (PayPalTxns, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	txnsJSON := PayPalTxnsJSON{}
	if err := dec.Decode(&txnsJSON); err != nil {
		return nil, err
	}

	return txnsJSON.Transactions, nil
}

func savePayPalTxnsToFile(filename string, txns PayPalTxns) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	txnsJSON := PayPalTxnsJSON{Transactions: txns}

	return enc.Encode(&txnsJSON)
}

func payPalTxnsFileName(year, month int) string {
	return fmt.Sprintf("%s/paypal-%d-%02d.json", dataDir, year, month)
}
