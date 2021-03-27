package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/leavengood/donation_tracker/paypal"
)

// Config is the main configuration for the whole program
type Config struct {
	PayPal *paypal.Config `json:"paypal"`

	// For getting the EUR to USD conversion rate
	FixerIoAccessKey string `json:"fixer_io_access_key"`

	// For updating the donations.json file on cdn.haiku-os.org
	Minio struct {
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
	} `json:"minio"`
}

func (c *Config) Validate() error {
	errorList := []string{}

	if c.PayPal.Endpoint == "" {
		errorList = append(errorList, "no PayPal endpoint was provided")
	}
	if c.PayPal.User == "" {
		errorList = append(errorList, "no PayPal user was provided")
	}
	if c.PayPal.Password == "" {
		errorList = append(errorList, "no PayPal password was provided")
	}
	if c.PayPal.Signature == "" {
		errorList = append(errorList, "no PayPal signature was provided")
	}

	if c.FixerIoAccessKey == "" {
		errorList = append(errorList, "no Fixer.io access key was provided")
	}

	if c.Minio.AccessKeyID == "" {
		errorList = append(errorList, "no Minio access key ID was provided")
	}
	if c.Minio.SecretAccessKey == "" {
		errorList = append(errorList, "no Minio secret access key was provided")
	}

	if len(errorList) == 0 {
		return nil
	}

	return errors.New(strings.Join(errorList, ", "))
}

const ConfigFile = "config.json"

var config Config

func LoadConfig() error {
	content, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return config.Validate()
}
