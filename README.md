Haiku Donation Tracker
======================

This is a small project used to automate the tracking of donations to [Haiku](https://www.haiku-os.org).
Currently donations are primarily received from PayPal, with occasional checks or larger donations
made directly to the [Haiku, Inc.](http://www.haiku-inc.org) bank account.

## To Build

```
go build
```

## To Run

```
./donation_tracker
```

This requires PayPal API credentials in config.json as described below.

## Code Organization

* `config.go`: Contains code for loading a simple `config.json` file containing PayPal API
credentials and other config information.

* `currency.go`: Contains the code for getting the EUR to USD conversion rate from "fixer.io".

* `main.go`: Contains the main function which orchestrates the overall process described above.

* `nvp.go`: Provides support for the unique "Name Value Pair" (NVP) format returned from PayPal API
calls.

* `paypal.go`: Provides all the code for handling, sorting, filtering and summarizing the PayPal
transactions.

* `paypal_api.go`: The code for calling the PayPal API which makes use of the NVP code.

* `summary.go`: All the code for overall summaries.

* `transactions.go`: Probably should be renamed. Loads the bank transactions from the CSV file described above.
