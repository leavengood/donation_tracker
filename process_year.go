package main

import (
	"fmt"
	"time"
)

// ProcessYear will take the provided year and EUR to USD conversion rate and
// perform the summary process which involves loading current data for the given
// year, getting any missing data, and then summarizing it all.
func ProcessYear(year int, eurToUsdRate float32) (*DonationSummary, error) {
	// TODO: Fill out new process using the ideal process text file.

	currentYear, currentMonth, _ := time.Now().UTC().Date()
	if year != currentYear {
		return nil, fmt.Errorf("Processing older years like %d is not currently supported", year)
	}

	// Otherwise assume we are processing an old year, and do all months

	// Load current files for the year
	fm, err := NewPayPalFileManager(year)
	if err != nil {
		return nil, err
	}

	// See what months are missing
	missing := fm.GetMissingMonths()
	// If the current year is being processed, remove future months
	if year == currentYear {
		newMissing := []int{}
		for _, month := range missing {
			if month > int(currentMonth) {
				break
			}
			newMissing = append(newMissing, month)
		}
		missing = newMissing
	}
	fmt.Printf("The missing months are: %v\n", missing)
	// t := fm.GetLatestTransaction()
	// fmt.Printf("The latest transaction is: %#v, with timestamp: %s\n", t, t.Timestamp)

	// Get missing months from PayPal API and save them
	var currentMonthTxns PayPalTxns
	for _, month := range missing {
		fmt.Printf("Fetching PayPal transactions for %02d/%d...\n", month, year)
		txns := FetchPayPalTxnsForMonth(year, month)
		fmt.Printf("    There are %d transactions\n", len(txns))

		// Save all but the current month in the current year
		if year == currentYear && month == int(currentMonth) {
			fmt.Println("    Not saving the current month")
			currentMonthTxns = txns
		} else {
			fmt.Println("    Saving this month")
			if err := fm.SaveMonth(month, txns); err != nil {
				return nil, err
			}
		}
	}

	summaries := MonthlySummaries{}
	summarizeMonth := func(month time.Month, txns PayPalTxns) {
		donations, other := txns.FilterDonations()
		monthStr := Colorize(34, fmt.Sprintf("%s %d", month, year))
		fmt.Printf("%s: There are %d donations from %d transactions\n",
			monthStr, len(donations), len(txns))
		sums := donations.Summarize()
		if len(sums) > 1 {
			fmt.Printf("    WARNING: multiple months found in summary for %s\n", monthStr)
		}
		monthSummary := sums[month]
		// At the beginning of the month in the current year, this could be empty
		if monthSummary != nil {
			fmt.Printf("    Donations: %s\n", monthSummary)
			summaries[month] = monthSummary
		}
		if len(other) > 0 {
			fmt.Println("    Other transactions:")
			for _, t := range other {
				fmt.Printf("        %s\n", t)
			}
		}
	}

	fmt.Printf("\n%s\n", Colorize(32, "Monthly Totals"))

	// Summarize each full month
	months := fm.GetExistingMonths()
	for _, month := range months {
		summarizeMonth(time.Month(month), fm.Months[month])
	}

	// Summarize current month transactions, if there are any
	if len(currentMonthTxns) > 0 {
		// Print them out first
		for _, txn := range currentMonthTxns {
			fmt.Println(txn)
		}
		summarizeMonth(currentMonth, currentMonthTxns)
	}

	// TODO: Add in special transactions to each monthly summary

	// Create totals and return the summary
	total := summaries.Total()
	fmt.Printf("\nTotal: %s\n", total)
	grossTotal := total.GrossTotal()
	fmt.Printf("Gross Total: %s\n", grossTotal)
	grandTotal := grossTotal.GrandTotal(eurToUsdRate)
	fmt.Println(Colorize(33, fmt.Sprintf("Grand Total (at EUR to USD rate of %f): %.02f",
		eurToUsdRate, grandTotal)))

	return &DonationSummary{
		UpdatedAt:      time.Now().UTC(),
		UsdDonations:   grossTotal["USD"],
		EurDonations:   grossTotal["EUR"],
		EurToUsdRate:   eurToUsdRate,
		TotalDonations: grandTotal,
	}, nil
}