package main

import (
	"fmt"
	"log"
)

func main() {
	data, err := CallPayPalNvpApi("TransactionSearch", "117.0",
		NameValues{"STARTDATE": "2014-09-01T00:00:00Z"})
	if err != nil {
		log.Fatal(err)
	}

	nvp := ParseNvpData(data)

	if nvp.Successful() {
		ppts := PayPalTxnsFromNvp(nvp)

		ppts.Sort()

		fmt.Printf("There are %v transactions:\n", len(ppts))
		for _, ppt := range ppts {
			fmt.Println(ppt)
		}

		donations := ppts.OnlyDonations()
		fmt.Printf("There are %v donations\n", len(donations))

		cc := donations.TotalByCurrency()
		for currency, amt := range cc {
			fmt.Printf("Total for %s: %.02f\n", currency, amt)
		}

		eurToUsdRate, err := GetExchangeRate()
		//eurToUsdRate := float32(1.2436)
		if err == nil {
			fmt.Printf("Grand Total (at EUR to USD rate of %f): %.02f\n", eurToUsdRate, cc.GrandTotal(eurToUsdRate))
		}
	} else {
		log.Fatalf("API call was not successful: %v\n", data)
	}
}
