package main

import (
	"fmt"
	"time"

	"github.com/leavengood/donation_tracker/other"
	"github.com/leavengood/donation_tracker/paypal"
	"github.com/leavengood/donation_tracker/util"
)

func AddTransactions(year int, m util.MonthlySummaries) {
	transactions, err := other.TransactionsFromCsv(fmt.Sprintf("transactions-%d.csv", year))
	if err != nil {
		fmt.Printf("Transaction file could not be found for %d, will assume there are no other transactions.\n", year)
		return
	}

	fmt.Printf("Adding %d transactions from CSV to the monthly summaries...\n", len(transactions))
	for _, t := range transactions {
		_, month, _ := t.Date.Date()
		fmt.Printf("    Merging in transaction [%s] to month %s\n", util.Colorize(util.Green, t.String()), month)
		m.ForMonth(month).AddOneTime(t.Amt, t.FeeAmt, t.CurrencyCode)
	}
}

func SummarizeYear(year int, eurToUsdRate float32, fm *paypal.FileManager) *DonationSummary {
	summaries := util.MonthlySummaries{}
	// TODO: Extract this so it can be used for the one month process. Maybe put it into
	// MonthlySummaries itself.
	summarizeMonth := func(month time.Month, txns paypal.Transactions) {
		donations, other := txns.FilterDonations()
		monthStr := util.Colorize(util.Blue, fmt.Sprintf("%s %d", month, year))
		fmt.Printf("%s: There are %d donations from %d transactions\n",
			monthStr, len(donations), len(txns))
		// TODO: Maybe this Summarize in PayPalTxns needs to be moved into
		// MonthlySummaries, so this all becomes less awkward.
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
			fmt.Println("\n    Other transactions:")
			for _, t := range other {
				fmt.Printf("        %s\n", t)
			}
		}
		fmt.Println("")
	}

	fmt.Printf("%s\n\n", util.Colorize(util.Green, "Monthly Totals"))

	// Summarize each month
	months := fm.GetExistingMonths()
	for _, month := range months {
		summarizeMonth(time.Month(month), fm.Months[month])
	}

	// Add in special transactions to each monthly summary
	AddTransactions(year, summaries)

	// Create totals and return the summary
	total := summaries.Total()
	fmt.Printf("\nTotal: %s\n", total)
	grossTotal := total.GrossTotal()
	fmt.Printf("Combined Total: %s\n", grossTotal)
	grandTotal := grossTotal.GrandTotal(eurToUsdRate)
	fmt.Println(util.Colorize(util.Yellow, fmt.Sprintf("Grand Total (at EUR to USD rate of %f): %.02f",
		eurToUsdRate, grandTotal)))

	return &DonationSummary{
		UpdatedAt:      time.Now().UTC(),
		UsdDonations:   grossTotal["USD"],
		EurDonations:   grossTotal["EUR"],
		EurToUsdRate:   eurToUsdRate,
		TotalDonations: grandTotal,
	}
}

// ProcessYear will take the provided year and EUR to USD conversion rate and
// perform the summary process which involves loading current data for the given
// year, getting any missing data, and then summarizing it all.
func ProcessYear(client *paypal.Client, year int, eurToUsdRate float32) (*DonationSummary, error) {
	currentYear, currentMonth, _ := time.Now().UTC().Date()

	// Load current files for the year
	fm, err := paypal.NewFileManager(year)
	if err != nil {
		return nil, err
	}

	// First deal with the latest month we have saved. It could be several
	// months ago depending on how long it has been between runs.
	latest := fm.GetLatestMonth()
	if latest != 0 {
		// Assume this is a partial month
		t := fm.GetLatestTransaction()
		fmt.Printf("The latest transaction is: %#v, with timestamp: %s\n", t, t.Timestamp)
		tYear, tMonth, tDay := t.Timestamp.Date()
		// Start from the beginning of this day so we don't miss anything
		startDate := fmt.Sprintf("%d-%02d-%02dT00:00:00Z", tYear, tMonth, tDay)
		fmt.Printf("Fetching PayPal transactions newer than: %s\n", startDate)
		newTxns := client.GetTransactions(startDate, paypal.GetEndDate(tYear, int(tMonth)))
		fmt.Printf("Found %d new transactions\n", len(newTxns))
		previous := fm.Months[latest]
		// Merge will remove any duplicates
		txns := previous.Merge(newTxns)

		// Only save if we got new transactions
		if len(txns) > len(previous) {
			if err := fm.SaveMonth(latest, txns); err != nil {
				return nil, err
			}
		} else {
			fmt.Printf("No new transactions found for latest month: %d\n", latest)
		}
	}

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

	// Get missing months from PayPal API and save them
	for _, month := range missing {
		if err := client.GetAndSaveMonth(year, month, fm); err != nil {
			return nil, err
		}
	}

	fmt.Println("")

	return SummarizeYear(year, eurToUsdRate, fm), nil
}
