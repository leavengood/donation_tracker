package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/leavengood/donation_tracker/paypal"
	"github.com/leavengood/donation_tracker/util"
)

// var month = flag.Int("month", 0, "Provide a month number for which you want to fetch fresh data")
// var year = flag.Int("year", 0, "Provide a year for which you want to fetch fresh data")
// var skipUpload = flag.Bool("skip-upload", false, "Skip uploading the donation summary to the CDN")

const usage = `Usage: %s [command]

If no command is given, the default command of "update" is performed for the
current year.

Commands:
    update [year] [--skip-upload | -s]
        Update the donation information for the given year, defaulting to the
        current year. Summary information is uploaded as JSON to the Haiku CDN
        for the current year only unless the --skip-upload or -s flag are
        provided.

    summarize [year]
        Provide a summary of a given year, defaulting to the current year. No
        new data is downloaded.

    fetch <month> [year]
        Fetch and save a single month of transactions from PayPal given a
        numeric month and optionally a year. The default year is the current
        year. Overwrites any existing data.

    help
        Show this usage.

A config file named config.json should be defined as described in the README.`

func main() {
	exit := func(msg string, exitCode int) {
		fmt.Println(msg)
		os.Exit(exitCode)
	}

	err := LoadConfig()
	if err != nil {
		exit(fmt.Sprintf("Could not load config file %v because of error: %v\n", ConfigFile, err), 1)
	}

	// Default command is update
	cmd := "update"
	args := os.Args
	if len(args) > 1 {
		cmd = args[1]
	}

	introPrint := func(msg string) {
		fmt.Printf("%s %s...\n\n", util.Colorize(util.Green, "✷"), msg)
	}

	greenCheck := util.Colorize(util.Green, "✓")
	blueArrow := util.Colorize(util.Blue, "❯")

	// currentYear, currentMonth, _ := time.Now().UTC().Date()
	currentYear, currentMonth, _ := time.Now().UTC().Date()

	// isCurrentYear := *year == 0 || *year == currentYear
	// fmt.Printf("Is current year: %v\n", isCurrentYear)

	getExchangeRate := func() float32 {
		fmt.Print("Fetching exchange rate for EUR to USD...")
		rate, err := GetExchangeRate(config.FixerIoAccessKey)
		if err != nil {
			// FIXME: Provide a helper for error printing and exit
			exit(fmt.Sprintf("Error: could not get EUR to USD exchange rate: %v", err), 1)
		}
		if rate == float32(0) {
			exit("Error: we got a 0 exchange rate from fixer.io, is the access key correct?", 1)
		}
		fmt.Printf("got rate of %f.\n\n", rate)
		return rate
	}

	getYear := func(yearStr string) int {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			exit(fmt.Sprintf("Error: invalid year provided: %s", yearStr), 1)
		}
		if year < 2010 || year > currentYear {
			exit(fmt.Sprintf("Error: Please provide a year between 2010 and %d", currentYear), 1)
		}

		return year
	}

	getMonth := func(monthStr string, year int) int {
		month, err := strconv.Atoi(monthStr)
		if err != nil {
			exit(fmt.Sprintf("Error: invalid month provided: %s", monthStr), 1)
		}
		maxMonth := 12
		if year == currentYear {
			maxMonth = int(currentMonth)
		}
		if month < 1 || month > maxMonth {
			exit(fmt.Sprintf("Error: Please provide a month between 1 and %d", maxMonth), 1)
		}

		return month
	}

	year := currentYear

	util.PrintLogo()

	client := paypal.NewClient(config.PayPal)

	switch cmd {
	case "help", "-help", "--help", "-h":
		exit(fmt.Sprintf(usage, args[0]), 0)

	case "update":
		skipUpload := false
		extraMsg := ""

		if len(args) == 3 {
			year = getYear(args[2])
		}

		if len(args) == 4 {
			flag := args[3]
			if flag == "--skip-upload" || flag == "-s" {
				skipUpload = true
				extraMsg = ", skipping upload of data."
			} else {
				fmt.Printf("Warning: unknown flag passed to update: %s\n", flag)
			}
		}

		introPrint(fmt.Sprintf("Updating donation information for %d%s", year, extraMsg))

		// Start with this so we fail fast if it has an error
		eurToUsdRate := getExchangeRate()

		ds, err := ProcessYear(client, year, eurToUsdRate)
		if err != nil {
			exit(fmt.Sprintf("Error: could not process year %d: %v\n", year, err), 1)
		}

		fmt.Printf("Donation Summary: %#v\n", ds)

		// Update the JSON file, if this is the current year and the skip flag was not set
		if year == currentYear && !skipUpload {
			fmt.Printf("%s Uploading donation summary...\n", blueArrow)
			err = UploadJson(ds)
			if err != nil {
				log.Fatalln(err)
			}
		}

		fmt.Printf("\n%s Update complete!\n", greenCheck)

	case "summarize":
		if len(args) == 3 {
			year = getYear(args[2])
		}

		fm, err := paypal.NewFileManager(year)
		if err != nil {
			exit(fmt.Sprintf("Error: could not load PayPal files for year %d: %v\n", year, err), 1)
		}
		latest := fm.GetLatestTransaction()
		if latest != nil {
			introPrint(fmt.Sprintf("The latest transaction is from %s", util.FormatDate(latest.Timestamp)))
		} else {
			introPrint(fmt.Sprintf("There does not seem to be any PayPal transactions for %d", year))
		}

		SummarizeYear(year, getExchangeRate(), fm)

	case "fetch":
		if len(args) < 3 {
			exit("Error: please provide a month to fetch", 1)
		}

		if len(args) == 4 {
			year = getYear(args[3])
		}

		month := getMonth(args[2], year)

		fm := paypal.NewEmptyFileManager(year)
		if err := client.GetAndSaveMonth(year, month, fm); err != nil {
			exit(fmt.Sprintf("Error: could not save transactions: %s", err), 1)
		}

	case "personal":
		if len(args) == 3 {
			year = getYear(args[2])
		}

		fm, err := paypal.NewFileManager(year)
		if err != nil {
			exit(fmt.Sprintf("Error: could not load PayPal files for year %d: %v\n", year, err), 1)
		}

		people := map[string]*util.Person{}
		for _, txns := range fm.Months {
			for _, t := range txns {
				if t.IsDonation() || t.IsSubscription() {
					// key := fmt.Sprintf("%s <%s>", t.Name, t.Email)
					key := t.Email
					person, found := people[key]
					if !found {
						person = &util.Person{
							Name:  t.Name,
							Email: t.Email,
							Total: util.CurrencyAmounts{},
							Count: 0,
						}
						people[key] = person
					}
					person.Total[t.CurrencyCode] += t.Amt
					person.Count++
				}
			}
		}

		fmt.Printf("There were donations from %d people:\n", len(people))
		for _, person := range people {
			fmt.Printf("  %s: %s (%d)\n",
				util.Colorize(util.Yellow, fmt.Sprintf("%s <%s>", person.Name, person.Email)),
				person.Total, person.Count)
		}
	default:
		fmt.Printf("Error: Unknown command %s.\n\n", cmd)
		exit(fmt.Sprintf(usage, args[0]), 1)
	}

	// TODO: Extract this into a function
	// if *month != 0 {
	// 	if *year == 0 {
	// 		*year = currentYear
	// 	}

	// 	if *year == currentYear && (*month < 1 || *month > int(currentMonth)) {
	// 		fmt.Printf("Please choose a month between 1 and %d\n", currentMonth)
	// 		os.Exit(1)
	// 	}

	// 	fmt.Printf("You chose month %d and year %d\n", *month, *year)
	// 	ppts := FetchPayPalTxnsForMonth(*year, *month)

	// 	fmt.Printf("There are %v transactions:\n", len(ppts))
	// 	for _, ppt := range ppts {
	// 		fmt.Println(ppt)
	// 	}

	// 	// Save to file
	// 	filename := payPalTxnsFileName(*year, *month)
	// 	fmt.Printf("Attempting to save %d transactions to file %s\n", len(ppts), filename)
	// 	err := savePayPalTxnsToFile(filename, ppts)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	donations, other := ppts.FilterDonations()
	// 	fmt.Printf("\n\nThere are %v donations:\n", len(donations))
	// 	for _, ppt := range donations {
	// 		fmt.Println(ppt)
	// 	}
	// 	fmt.Printf("\n\nThere are %v other transactions:\n", len(other))
	// 	for _, ppt := range other {
	// 		fmt.Println(ppt)
	// 	}

	// cc := donations.TotalByCurrency()
	// for currency, amt := range cc {
	// 	fmt.Printf("Total for %s: %.02f\n", currency, amt)
	// }

	// eurToUsdRate, err := GetExchangeRate()
	// //eurToUsdRate := float32(1.2436)
	// if err == nil {
	// 	fmt.Printf("Grand Total (at EUR to USD rate of %f): %.02f\n", eurToUsdRate, cc.GrandTotal(eurToUsdRate))
	// }

	// 	currentSummaries := donations.Summarize()
	// 	for month, summary := range currentSummaries {
	// 		fmt.Printf("%s: %v\n", month, summary)
	// 		fmt.Printf("Gross Total: %v\n", summary.GrossTotal())
	// 		fmt.Printf("Net Total: %v\n\n", summary.NetTotal())
	// 		b, _ := json.Marshal(summary)
	// 		fmt.Println(string(b))
	// 	}
	// 	os.Exit(0)
	// }
	// fmt.Println("Did not get a month")
	// summaryProcess()
	// fm, err := NewPayPalFileManager(2019)
	// // months, err := loadPayPalFiles("data", 2019)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Found %d months worth of files\n", len(fm.Months))
	// for month, txns := range fm.Months {
	// 	fmt.Printf("There are %d transactions for month %d\n", len(txns), month)
	// }
	// fmt.Printf("The most recent month is: %d\n", fm.GetLatestMonth())
	// fmt.Printf("The missing months are: %v\n", fm.GetMissingMonths())

	// Start with this so we fail fast if it has an error
	// eurToUsdRate, err := GetExchangeRate(config.FixerIoAccessKey)
	// if err != nil {
	// 	log.Fatalf("could not get EUR to USD exchange rate because of error: %v\n", err)
	// }
	// if eurToUsdRate == float32(0) {
	// 	log.Fatalln("we got a 0 exchange rate from fixer.io, is the access key correct?")
	// }

	// ds, err := ProcessYear(currentYear, eurToUsdRate)
	// if err != nil {
	// 	log.Fatalf("could not process year %d due to error: %v\n", currentYear, err)
	// }

	// fmt.Printf("Donation Summary: %#v\n", ds)

	// // Update the JSON file, if this is the current year and the skip flag was not set
	// if isCurrentYear && !*skipUpload {
	// 	fmt.Println("Uploading donation summary...")
	// 	err = UploadJson(ds)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	fmt.Println("Skipping upload of donation summary")
	// }
}
