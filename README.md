# Haiku Donation Tracker

This is a small project used to automate the tracking of donations to [Haiku](https://www.haiku-os.org).
Currently donations are primarily received from PayPal, with occasional checks or larger donations
made directly to the [Haiku, Inc.](http://www.haiku-inc.org) bank account.

This tool has multiple small commands it can perform against the donation information it collects.

## Commands

### `update`

The default and most important command is `update`, which performs the standard process to update
the donation meter. This calls the PayPal NVP (Name Value Pair) API to get a list of transactions.
Those are fetched by the month and saved into JSON files in the data directory. Tthe most recently
updated month could have a partial list of transactions. The update process determines the most
recently updated month and ensures that no transactions are missed when fetching new ones.

These saved PayPal transactions are filtered and grouped into one-time and subscription donations and
then totaled based on currency (currently just USD and EUR.) The current EUR to USD exchange rate is
fetched from the "fixer.io" API and used to convert the EUR donation total into USD to make a grand
total. That information is then saved into a `donation.json` file which is uploaded to https://cdn.haiku-os.org/haiku-inc.

The reason transactions are grouped by type of donation (one-time and subscription) is that information
is intended to be used to update a monthly summary of donations, but that is not done yet.

Information about subscription creation and cancellation is also received from PayPal and is
printed during the `update` process. This may also be included in the future monthly donation report.

TODO: Document other commands

## To Build

```
go build
```

## To Run

```
./donation_tracker <command>
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
