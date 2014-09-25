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
		for _, ppt := range(ppts) {
			fmt.Println(ppt)
		}
	} else {
		log.Fatalf("API call was not successful: %v\n", data)
	}
}
