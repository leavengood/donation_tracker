package paypal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const dataDir = "data"

// FileManager manages files containing PayPal transactions fetched from
// the PayPal API.
type FileManager struct {
	Year   int
	Months map[int]Transactions
}

// NewFileManager will load any PayPal files for the given year from the
// data directory, and can be used to save new transactions there.
func NewFileManager(year int) (*FileManager, error) {
	months, err := loadPayPalFiles(dataDir, year)
	if err != nil {
		return nil, err
	}

	return &FileManager{
		Year:   year,
		Months: months,
	}, nil
}

func NewEmptyFileManager(year int) *FileManager {
	return &FileManager{
		Year:   year,
		Months: map[int]Transactions{},
	}
}

// GetLatestMonth returns the latest month with transactions loaded by this
// file manager. When managing the current year, this helps determine what new
// data needs to be fetched from the PayPal API.
func (p *FileManager) GetLatestMonth() int {
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
func (p *FileManager) GetLatestTransaction() *Transaction {
	txns := p.Months[p.GetLatestMonth()]

	// The transactions should be sorted with the latest last
	if len(txns) > 0 {
		return txns[len(txns)-1]
	}

	return nil
}

// GetExistingMonths will return all months which are currently stored in this
// file manager.
func (p *FileManager) GetExistingMonths() []int {
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
func (p *FileManager) GetMissingMonths() []int {
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
func (p *FileManager) SaveMonth(month int, txns Transactions) error {
	filename := payPalTxnsFileName(p.Year, month)
	if err := savePayPalTxnsToFile(filename, txns); err != nil {
		return err
	}

	p.Months[month] = txns

	return nil
}

func loadPayPalFiles(baseDir string, year int) (map[int]Transactions, error) {
	result := map[int]Transactions{}

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

// TransactionsJSON is the structure that will be saved to and loaded from files.
type TransactionsJSON struct {
	Transactions Transactions `json:"transactions"`
}

func loadPayPalTxnsFromFile(filename string) (Transactions, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(f)
	txnsJSON := TransactionsJSON{}
	if err := dec.Decode(&txnsJSON); err != nil {
		return nil, err
	}

	return txnsJSON.Transactions, nil
}

func savePayPalTxnsToFile(filename string, txns Transactions) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	txnsJSON := TransactionsJSON{Transactions: txns}

	return enc.Encode(&txnsJSON)
}

func payPalTxnsFileName(year, month int) string {
	return fmt.Sprintf("%s/paypal-%d-%02d.json", dataDir, year, month)
}
