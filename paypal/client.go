package paypal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/leavengood/donation_tracker/util"
)

const (
	// The API endpoint to call, provided in the config file
	EndpointKey = "ENDPOINT"

	// Fields needed to make a request
	MethodKey  = "METHOD"
	VersionKey = "VERSION"

	// These are security params, provided in the config file
	UserKey      = "USER"
	PasswordKey  = "PWD"
	SignatureKey = "SIGNATURE"
)

type Client struct {
	config *Config
	Debug  bool
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) GetTransactions(startDate, endDate string) Transactions {
	fmt.Printf("Getting PayPal data from %s to %s\n", startDate, endDate)

	nv := NameValues{"STARTDATE": startDate}
	if endDate != "" {
		nv["ENDDATE"] = endDate
	}
	data, err := callPayPalNvpApi(c.config, "TransactionSearch", "117.0", nv)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Make this easier to use for debugging
	// ioutil.WriteFile("paypal.nvp", []byte(data), 0700)

	// fileData, _ := ioutil.ReadFile("paypal.nvp")
	// data := string(fileData)

	nvp := ParseNvpData(data)

	if !nvp.Successful() {
		// TODO: Return error
		log.Fatalf("API call was not successful: %v\n", data)
	}

	ppts := TransactionsFromNvp(nvp)
	ppts.Sort()

	return ppts
}

// TODO: Maybe move...
func GetEndDate(year, month int) string {
	// if month == 12 {
	// 	year++
	// 	month = 0
	// }
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	_, _, day := lastOfMonth.Date()
	return fmt.Sprintf("%d-%02d-%02dT23:59:59Z", year, month, day)
}

func (c *Client) GetTransactionsForMonth(year, month int) Transactions {
	startDate := fmt.Sprintf("%d-%02d-01T00:00:00Z", year, month)

	// TODO: better error handling
	return c.GetTransactions(startDate, GetEndDate(year, month))
}

func (c *Client) GetAndSaveMonth(year, month int, fm *FileManager) error {
	monthStr := util.Colorize(util.Green, fmt.Sprintf("%s %d", time.Month(month), year))
	fmt.Printf("Fetching PayPal transactions for %s...", monthStr)
	txns := c.GetTransactionsForMonth(year, month)
	fmt.Printf("there are %d transactions, saving to JSON.\n", len(txns))
	return fm.SaveMonth(month, txns)
}

func callPayPalNvpApi(config *Config, method string, version string, params NameValues) (string, error) {
	v := url.Values{}
	v.Set(MethodKey, method)
	v.Set(VersionKey, version)
	v.Set(UserKey, config.User)
	v.Set(PasswordKey, config.Password)
	v.Set(SignatureKey, config.Signature)

	for name, value := range params {
		v.Set(name, value)
	}

	resp, err := http.PostForm(config.Endpoint, v)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
