package main

import (
	"fmt"
	"log"
)

func main() {
	data, err := CallPayPalNvpApi("TransactionSearch", "117.0",
		NameValues{"STARTDATE": "2014-09-25T09:00:00Z"})
	if err != nil {
		log.Fatal(err)
	}
	nvp := ParseNvpData(data)

	if nvp.Successful() {
		ppts := PayPalTxnsFromNvp(nvp)

		fmt.Println(ppts[0])
	}
}
