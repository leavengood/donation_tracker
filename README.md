Haiku Donation Tracker
======================

This is a small project used to automate the tracking of donations to [Haiku](https://www.haiku-os.org).
Currently donations are primarily received from PayPal, with occasional checks or larger donations
made directly to the [Haiku, Inc.](http://www.haiku-inc.org) bank account.

This program calls the PayPal NVP (Name Value Pair) API to get a list of transactions. It then
filters and groups those into one-time and subscription donations and totals then based on currency
(currently just USD and EUR.) The current EUR to USD exchange rate is gotten from the "fixer.io" API
and used to convert the EUR donation total into USD to make a grand total. That information is then
saved into a `donation.json` file which is uploaded to https://cdn.haiku-os.org/haiku-inc.

Then if a complete months worth of transactions was received it is saved into a running summary file
by the month, so we don't have to get that months transactions again.

Transaction detailed information is currently not kept, but that might be a nice addition.

Information about subscription creation and cancellation is also received from PayPal but is not
currently used.

## To Build

```
go build
```

## To Run

```
./donation_tracker
```

This requires various API credentials in config.json in this format:

```
{
  "paypal": {
    "user": "",
    "password": "",
    "signature": "",
    "endpoint": ""
  },
  "fixer_io_access_key": "",
  "minio": {
    "access_key_id": "",
    "secret_access_key": ""
  }
}
```

The PayPal credentials are to get the transactions. The "fixer.io" access key is
for getting the EUR to USD conversion rate. The Minio credentials are for uploading
a JSON file with the donation summary information to https://cdn.haiku-os.org.

If any config values are missing the code will not run.

## Code Organization

* `config.go`: Contains code for loading a simple `config.json` file containing PayPal API
credentials and other config information.

* `currency.go`: Contains the code for getting the EUR to USD conversion rate from "fixer.io".

* `json_upload.go`: Contains the code for uploading the summary JSON file to the Haiku Minio server.

* `main.go`: Contains the main function which orchestrates the overall process described above.

* `nvp.go`: Provides support for the unique "Name Value Pair" (NVP) format returned from PayPal API
calls.

* `paypal.go`: Provides all the code for handling, sorting, filtering and summarizing the PayPal
transactions.

* `paypal_api.go`: The code for calling the PayPal API which makes use of the NVP code.

* `summary.go`: All the code for overall summaries.

* `transactions.go`: Probably should be renamed. Loads the bank transactions from the CSV file described above.
