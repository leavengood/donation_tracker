package paypal

// Config has various config values for getting PayPal transations with their NVP API
type Config struct {
	Endpoint  string `json:"endpoint"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Signature string `json:"signature"`
}
