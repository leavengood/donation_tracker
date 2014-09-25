package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	SandboxEndPoint = "https://api-3t.sandbox.paypal.com/nvp"

	UserKey = "USER"
	PasswordKey = "PWD"
	SignatureKey = "SIGNATURE"
	MethodKey = "METHOD"
	VersionKey = "VERSION"

	User = "jb-us-seller_api1.paypal.com"
	Password = "WX4WTU3S8MY44S7F"
	Signature = "AFcWxV21C7fd0v3bYYYRCpSSRl31A7yDhhsPUU2XhtMoZXsWHFxu-RWy"
)

// func CallPayPalNvpApi(method string, version string, params NameValues) {
func CallPayPalNvpApi(method string, version string) (string, error) {
	v := url.Values{}
	v.Set(MethodKey, method)
	v.Set(VersionKey, version)
	v.Set(UserKey, User)
	v.Set(PasswordKey, Password)
	v.Set(SignatureKey, Signature)
	v.Set("STARTDATE", "2014-09-25T09:00:00Z")

	resp, err := http.PostForm(SandboxEndPoint, v)
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
