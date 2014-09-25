package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// The API endpoint to call, provided in the config file
	EndpointKey = "ENDPOINT"

	// Fields needed to make a request
	MethodKey = "METHOD"
	VersionKey = "VERSION"
	// These are security params, provided in the config file
	UserKey = "USER"
	PasswordKey = "PWD"
	SignatureKey = "SIGNATURE"
)

func CallPayPalNvpApi(method string, version string, params NameValues) (string, error) {
	v := url.Values{}
	v.Set(MethodKey, method)
	v.Set(VersionKey, version)
	v.Set(UserKey, config[UserKey])
	v.Set(PasswordKey, config[PasswordKey])
	v.Set(SignatureKey, config[SignatureKey])

	for name, value := range(params) {
		v.Set(name, value)
	}

	resp, err := http.PostForm(config[EndpointKey], v)
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
