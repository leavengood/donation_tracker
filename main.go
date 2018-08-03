package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func FetchPayPalTxns(startDate, endDate string) PayPalTxns {
	fmt.Printf("Getting PayPal data from %s\n", startDate)
	// startDate = "2015-06-01T00:00:00Z"
	// endDate := "2015-07-01T00:00:00Z"
	nv := NameValues{"STARTDATE": startDate}
	if endDate != "" {
		nv["ENDDATE"] = endDate
	}
	// NameValues{"STARTDATE": startDate, "ENDDATE": endDate}
	data, err := CallPayPalNvpApi("TransactionSearch", "117.0", nv)
	if err != nil {
		log.Fatal(err)
	}
	// ioutil.WriteFile("paypal.nvp", []byte(data), 0700)

	// fileData, _ := ioutil.ReadFile("paypal.nvp")
	// data := string(fileData)

	nvp := ParseNvpData(data)

	if !nvp.Successful() {
		log.Fatalf("API call was not successful: %v\n", data)
	}

	ppts := PayPalTxnsFromNvp(nvp)
	ppts.Sort()

	return ppts
}

func FetchPayPalTxnsForMonth(year, month int) PayPalTxns {
	startDate := fmt.Sprintf("%d-%02d-01T00:00:00Z", year, month)
	if month == 12 {
		year += 1
		month = 0
	}
	endDate := fmt.Sprintf("%d-%02d-01T00:00:00Z", year, month+1)

	return FetchPayPalTxns(startDate, endDate)
}

// TODO:
// - Make a command to get totals for each month fresh from PayPal.
// - Make a command to show a summary.
// - Create a String() method for PayPal transations so they look nicer in the shell.
// - Make small command-line program to update the donation meter in the Haiku
//   website database based on command-line arguments. This can be used from ssh.
// - Add code to create the donation_meter.txt in the haiku-inc.org site.
// - Make the summarizing more robust, so as not to accidentally increase when
//   we shouldn't (like what happened with June.)
// - Create monthly summary tables.
// - It would be nice to track subscription creation and cancellation.
// - When updating daily, we can show how many donations came in that day.

func summaryProcess() {
	// Start with this so we fail fast if it has an error
	eurToUsdRate, err := GetExchangeRate(config.FixerIoAccessKey)
	if err != nil {
		log.Fatalf("could not get EUR to USD exchange rate because of error: %v\n", err)
	}
	if eurToUsdRate == float32(0) {
		log.Fatalln("we got a 0 exchange rate from fixer.io, is the access key correct?")
	}

	year, currentMonth, _ := time.Now().Date()

	// Load previous summaries
	summariesPrefix := fmt.Sprintf("summaries-%d", year)
	previousSummaries, err := LoadSummaries(summariesPrefix)
	if err != nil {
		fmt.Printf("Summary file could not be found for %d, will create if needed.\n", year)
	}

	// Load other transactions
	transactions, err := TransactionsFromCsv(fmt.Sprintf("transactions-%d.csv", year))
	if err != nil {
		fmt.Printf("Transaction file could not be found for %d, will assume there are no other transactions.\n", year)
	}

	// Get the start date for PayPal
	currentMonthInt := int(currentMonth)
	currentMonthToFetch := 1
	if len(previousSummaries) > 0 {
		// Start with the month after the latest month in the summaries
		currentMonthToFetch = int(previousSummaries.LatestMonth() + 1)
	}

	var ppts PayPalTxns
	if currentMonthInt-currentMonthToFetch > 2 {
		fmt.Println("We are pretty behind, fetching transactions month by month...")
		// We are pretty far behind, go month by month
		for i := currentMonthToFetch; i < currentMonthInt; i++ {
			ppts = append(ppts, FetchPayPalTxnsForMonth(year, i)...)
		}
	} else {
		startDate := fmt.Sprintf("%d-%02d-01T00:00:00Z", year, currentMonthToFetch)
		ppts = FetchPayPalTxns(startDate, "")
	}

	fmt.Printf("There are %v transactions:\n", len(ppts))
	for _, ppt := range ppts {
		fmt.Println(ppt)
	}

	donations, _ := ppts.FilterDonations()
	fmt.Printf("There are %v donations\n", len(donations))

	// cc := donations.TotalByCurrency()
	// for currency, amt := range cc {
	// 	fmt.Printf("Total for %s: %.02f\n", currency, amt)
	// }

	currentSummaries := donations.Summarize()
	// currentWithExtraTrans := make(MonthlySummaries)

	for month, summary := range currentSummaries {
		fmt.Printf("%s: %v\n", month, summary)
		fmt.Printf("Gross Total: %v\n", summary.GrossTotal())
		fmt.Printf("Net Total: %v\n\n", summary.NetTotal())
		// sCopy := *summary
		// currentWithExtraTrans[month] = &sCopy
	}

	// Merge in other transactions for current months
	for _, t := range transactions {
		_, month, _ := t.Date.Date()
		if summary, ok := currentSummaries[month]; ok {
			fmt.Printf("Merging in transaction %v\n", t)
			summary.AddOneTime(t.Amt, t.FeeAmt, t.CurrencyCode)
		}
	}

	// for _, t := range transactions {
	// 	fmt.Printf("Merging in transaction %v\n", t)
	// 	_, month, _ := t.Date.Date()
	// 	if summary, ok := previousSummaries[month]; ok {
	// 		fmt.Printf("Summary for month %d before: %v\n", month, summary)
	// 		summary.AddOneTime(t.Amt, t.FeeAmt, t.CurrencyCode)
	// 		fmt.Printf("Summary for month %d after: %v\n\n", month, summary)
	// 	} else if summary, ok := currentWithExtraTrans[month]; ok {
	// 		// TODO: Remove duplication!
	// 		fmt.Printf("Summary for month %d before: %v\n", month, summary)
	// 		summary.AddOneTime(t.Amt, t.FeeAmt, t.CurrencyCode)
	// 		fmt.Printf("Summary for month %d after: %v\n\n", month, summary)
	// 	}
	// }

	fmt.Printf("Previous: %v, Gross Total: %s\n", previousSummaries.Total(), previousSummaries.Total().GrossTotal())
	fmt.Printf("Current: %v, Gross Total: %s\n", currentSummaries.Total(), currentSummaries.Total().GrossTotal())
	total := previousSummaries.Total().GrossTotal().Add(currentSummaries.Total().GrossTotal())
	// cc := donations.TotalByCurrency()
	for currency, amt := range total {
		fmt.Printf("Total for %s: %.02f\n", currency, amt)
	}

	grandTotal := total.GrandTotal(eurToUsdRate)
	fmt.Printf("Grand Total (at EUR to USD rate of %f): %.02f\n", eurToUsdRate, grandTotal)
	fmt.Printf("Percentage pixels (from 128): %.02f\n", grandTotal/10000.0*128)

	// Update the JSON file
	err = UploadJson(&DonationSummary{
		UpdatedAt:      time.Now().UTC(),
		UsdDonations:   total["USD"],
		EurDonations:   total["EUR"],
		EurToUsdRate:   eurToUsdRate,
		TotalDonations: grandTotal,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Do we have more than one month summarized?
	if len(currentSummaries) > 1 {
		// FIXME: We save the separate transactions when we should just save the PayPal ones...
		fmt.Printf("Saving summaries for %d months\n", len(currentSummaries)-1)
		// Then let's merge it with the existing summaries and save to file
		for month, summary := range currentSummaries {
			// Don't save the current month
			if month != currentMonth {
				fmt.Printf("Saving summary for month %s: %s\n", month, summary)
				previousSummaries[month] = summary
			}
		}
		previousSummaries.Save(summariesPrefix)
	}
}

var month = flag.Int("month", 0, "Provide a month number for which you want to fetch fresh data")
var year = flag.Int("year", 0, "Provide a year for which you want to fetch fresh data")

func main() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf("could not load config file %v because of error: %v\n", ConfigFile, err)
	}

	flag.Parse()

	currentYear, currentMonth, _ := time.Now().Date()

	if *month != 0 {
		if *year == 0 {
			*year = currentYear
		}

		if *year == currentYear && (*month < 1 || *month > int(currentMonth)) {
			fmt.Printf("Please choose a month between 1 and %d\n", currentMonth)
			os.Exit(1)
		}

		fmt.Printf("You chose month %d and year %d\n", *month, *year)
		ppts := FetchPayPalTxnsForMonth(*year, *month)

		fmt.Printf("There are %v transactions:\n", len(ppts))
		for _, ppt := range ppts {
			fmt.Println(ppt)
		}

		donations, other := ppts.FilterDonations()
		fmt.Printf("\n\nThere are %v donations:\n", len(donations))
		for _, ppt := range donations {
			fmt.Println(ppt)
		}
		fmt.Printf("\n\nThere are %v other transactions:\n", len(other))
		for _, ppt := range other {
			fmt.Println(ppt)
		}

		// cc := donations.TotalByCurrency()
		// for currency, amt := range cc {
		// 	fmt.Printf("Total for %s: %.02f\n", currency, amt)
		// }

		// eurToUsdRate, err := GetExchangeRate()
		// //eurToUsdRate := float32(1.2436)
		// if err == nil {
		// 	fmt.Printf("Grand Total (at EUR to USD rate of %f): %.02f\n", eurToUsdRate, cc.GrandTotal(eurToUsdRate))
		// }

		currentSummaries := donations.Summarize()
		for month, summary := range currentSummaries {
			fmt.Printf("%s: %v\n", month, summary)
			fmt.Printf("Gross Total: %v\n", summary.GrossTotal())
			fmt.Printf("Net Total: %v\n\n", summary.NetTotal())
		}
		os.Exit(0)
	}
	// fmt.Println("Did not get a month")
	summaryProcess()
}
