package util

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"
)

type Donor struct {
	Name      string
	Email     string
	Total     CurrencyAmounts
	Count     int
	Anonymous bool
}

// TODO: Implement String()

type Donors []*Donor

func (p Donors) Len() int      { return len(p) }
func (p Donors) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type ByTotal struct{ Donors }

func (s ByTotal) Less(i, j int) bool {
	totalI := s.Donors[i].Total.GrandTotal(1)
	totalJ := s.Donors[j].Total.GrandTotal(1)
	// We want descending order
	return totalI >= totalJ
}

func (p Donors) Sort() {
	sort.Sort(ByTotal{p})
}

const DonorConfigFile = "donors.json"

type DonorConfig struct {
	Anonymous          []string          `json:"anonymous,omitempty"`
	EmailToCorrectName map[string]string `json:"names,omitempty"`

	anonMap map[string]bool
}

func (d *DonorConfig) Handle(p *Donor) {
	if d.anonMap == nil {
		d.anonMap = make(map[string]bool, len(d.Anonymous))
		for _, e := range d.Anonymous {
			d.anonMap[e] = true
		}
	}

	email := strings.ToLower(p.Email)
	if name, found := d.EmailToCorrectName[email]; found {
		p.Name = name
	}
	if d.anonMap[email] {
		p.Anonymous = true
	}
}

func LoadDonorConfig() (*DonorConfig, error) {
	content, err := ioutil.ReadFile(DonorConfigFile)
	if err != nil {
		return nil, err
	}

	var config DonorConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
