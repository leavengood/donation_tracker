package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/leavengood/donation_tracker/paypal"
	"github.com/leavengood/donation_tracker/util"
)

const usage = `
Usage: %s [command]

If no command is given, the default command of "update" is performed for the
current year.

Commands:
    update [-year int] [-skip-upload]
        Update the donation information for the given year, defaulting to the
        current year. Summary information is uploaded as JSON to the Haiku CDN
        for the current year only unless the --skip-upload flag is provided.

    summarize [-year int]
        Provide a summary of a given year, defaulting to the current year. No
        new data is downloaded.

    fetch <-month int> [-year int]
        Fetch and save a single month of transactions from PayPal given a
        numeric month and optionally a year. The default year is the current
        year. Overwrites any existing data.

    donors [-year int]
        Collect information for donors in the given year, defaulting to the
        current year, and print out their name, email, how much and how many
        times they have donated, in descending order of donation amount.

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
	exe := os.Args[0]
	cmd := "update"
	args := os.Args[1:]
	if len(args) > 0 {
		// This could be a flag for the default command of update
		if args[0][0] != '-' {
			cmd = args[0]
			args = args[1:]
		}
	}

	currentYear, currentMonth, _ := time.Now().UTC().Date()

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	year := currentYear
	flagSet.IntVar(&year, "year", currentYear, "Specifies the year to operate on")
	month := int(currentMonth)
	flagSet.IntVar(&month, "month", int(currentMonth), "Specifies the month to operate on")
	skipUpload := flagSet.Bool("skip-upload", false, "Skip uploading data to the server on the 'update' command")
	emails := flagSet.Bool("emails", false, "Print only emails in the 'donors' command")

	printUsage := func() {
		fmt.Println(fmt.Sprintf(usage, exe))
		fmt.Println("\nSupported flags:")
		flagSet.PrintDefaults()
	}
	flagSet.Usage = printUsage
	flagSet.Parse(args)

	introPrint := func(msg string) {
		fmt.Printf("%s %s...\n\n", util.Colorize(util.Green, "✷"), msg)
	}

	greenCheck := util.Colorize(util.Green, "✓")
	blueArrow := util.Colorize(util.Blue, "❯")

	getExchangeRate := func() float32 {
		fmt.Print("Fetching exchange rate for EUR to USD...")
		rate, err := GetExchangeRate(config.FixerIoAccessKey)
		if err != nil {
			exit(fmt.Sprintf("Error: could not get EUR to USD exchange rate: %v", err), 1)
		}
		if rate == float32(0) {
			exit("Error: we got a 0 exchange rate from fixer.io, is the access key correct?", 1)
		}
		fmt.Printf("got rate of %f.\n\n", rate)
		return rate
	}

	// Sanity check the year
	if year < 2010 || year > currentYear {
		exit(fmt.Sprintf("Error: Please provide a year between 2010 and %d", currentYear), 1)
	}

	util.PrintLogo()

	client := paypal.NewClient(config.PayPal)

	switch cmd {
	case "help":
		printUsage()
		os.Exit(0)

	case "update":
		extraMsg := ""
		if *skipUpload {
			extraMsg = ", skipping upload of data."
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
		if year == currentYear && !*skipUpload {
			fmt.Printf("%s Uploading donation summary...\n", blueArrow)
			err = UploadJson(ds)
			if err != nil {
				log.Fatalln(err)
			}
		}

		fmt.Printf("\n%s Update complete!\n", greenCheck)

	case "summarize":
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
		// Sanity check the month
		maxMonth := 12
		if year == currentYear {
			maxMonth = int(currentMonth)
		}
		if month < 1 || month > maxMonth {
			exit(fmt.Sprintf("Error: Please provide a month between 1 and %d", maxMonth), 1)
		}

		fm := paypal.NewEmptyFileManager(year)
		if err := client.GetAndSaveMonth(year, month, fm); err != nil {
			exit(fmt.Sprintf("Error: could not save transactions: %s", err), 1)
		}

	case "donors":
		fm, err := paypal.NewFileManager(year)
		if err != nil {
			exit(fmt.Sprintf("Error: could not load PayPal files for year %d: %v\n", year, err), 1)
		}
		config, err := util.LoadDonorConfig()
		if err != nil {
			exit(fmt.Sprintf("Error: could not load the donor config file: %v\n", err), 1)
		}

		donorMap := map[string]*util.Donor{}
		for _, txns := range fm.Months {
			for _, t := range txns {
				if t.IsDonation() || t.IsSubscription() {
					key := t.Email
					donor, found := donorMap[key]
					if !found {
						donor = &util.Donor{
							Name:  t.Name,
							Email: t.Email,
							Total: util.CurrencyAmounts{},
							Count: 0,
						}
						// Correct their name or set the anoymous flag
						config.Handle(donor)
						donorMap[key] = donor
					}
					donor.Total[t.CurrencyCode] += t.Amt
					donor.Count++
				}
			}
		}

		donors := make(util.Donors, 0, len(donorMap))
		for _, person := range donorMap {
			donors = append(donors, person)
		}
		donors.Sort()

		fmt.Printf("There were donations from %d donors:\n", len(donors))
		for _, person := range donors {
			if *emails {
				fmt.Println(person.Email)
			} else {
				anon := ""
				if person.Anonymous {
					anon = util.Colorize(util.BrightYellow, "{Wishes to be Anonymous}")
				}
				fmt.Printf("  %s: %s (%d) %s\n",
					util.Colorize(util.Yellow, fmt.Sprintf("%s <%s>", person.Name, person.Email)),
					person.Total, person.Count, anon)
			}
		}

	default:
		fmt.Printf("Error: Unknown command %s.\n\n", cmd)
		exit(fmt.Sprintf(usage, args[0]), 1)
	}
}
