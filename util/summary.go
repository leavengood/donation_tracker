package util

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Person struct {
	Name  string
	Email string
	Total CurrencyAmounts
	Count int
}

type Summary struct {
	OneTimeAmt        CurrencyAmounts
	OneTimeCount      int
	SubscriptionAmt   CurrencyAmounts
	SubscriptionCount int
	FeeAmt            CurrencyAmounts
}

func NewSummary() *Summary {
	result := new(Summary)

	result.OneTimeAmt = make(CurrencyAmounts)
	result.SubscriptionAmt = make(CurrencyAmounts)
	result.FeeAmt = make(CurrencyAmounts)

	return result
}

// func (s1 *Summary) Add(s2 *Summary) *Summary {

// result := NewSummary()

// for _, summary := range ms {
//   result.OneTimeAmt.Add(summary.OneTimeAmt)
//   result.OneTimeCount += summary.OneTimeCount
//   result.SubscriptionAmt.Add(summary.SubscriptionAmt)
//   result.SubscriptionCount += summary.SubscriptionCount
//   result.FeeAmt.Add(summary.FeeAmt)
// }

// return result
// }

func (s *Summary) String() string {
	return fmt.Sprintf("OneTime: %s (%d), Subscriptions: %s (%d), Fees: %s",
		s.OneTimeAmt, s.OneTimeCount, s.SubscriptionAmt, s.SubscriptionCount, s.FeeAmt)
}

func (s *Summary) AddOneTime(amt, fee float32, currency string) {
	s.OneTimeCount += 1
	s.OneTimeAmt[currency] += amt
	s.FeeAmt[currency] += fee
}

func (s *Summary) AddSubscription(amt, fee float32, currency string) {
	s.SubscriptionCount += 1
	s.SubscriptionAmt[currency] += amt
	s.FeeAmt[currency] += fee
}

func (s *Summary) GrossTotal() CurrencyAmounts {
	return s.OneTimeAmt.Add(s.SubscriptionAmt)
}

func (s *Summary) NetTotal() CurrencyAmounts {
	// The fees are negative, so they are added
	return s.GrossTotal().Add(s.FeeAmt)
}

type MonthlySummaries map[time.Month]*Summary

func (ms MonthlySummaries) ForMonth(month time.Month) *Summary {
	summary := ms[month]

	if summary == nil {
		summary = NewSummary()
		ms[month] = summary
	}

	return summary
}

func (ms MonthlySummaries) LatestMonth() time.Month {
	result := time.January

	if len(ms) > 0 {
		for m := range ms {
			if m > result {
				result = m
			}
		}
	}

	return result
}

func (ms MonthlySummaries) ByMonth(cb func(time.Month, *Summary)) {
	latest := ms.LatestMonth()
	for i := time.January; i <= latest; i++ {
		cb(i, ms[i])
	}
}

func (ms MonthlySummaries) Total() *Summary {
	result := NewSummary()

	for _, summary := range ms {
		// FIXME: This sucks
		result.OneTimeAmt = result.OneTimeAmt.Add(summary.OneTimeAmt)
		result.OneTimeCount += summary.OneTimeCount
		result.SubscriptionAmt = result.SubscriptionAmt.Add(summary.SubscriptionAmt)
		result.SubscriptionCount += summary.SubscriptionCount
		result.FeeAmt = result.FeeAmt.Add(summary.FeeAmt)
	}

	return result
}

func (ms MonthlySummaries) MarshalJSON() ([]byte, error) {
	// Convert to string keys
	stringKeys := make(map[string]*Summary)
	for m, s := range ms {
		stringKeys[m.String()] = s
	}

	return json.Marshal(stringKeys)
}

func (ms MonthlySummaries) UnmarshalJSON(data []byte) error {
	stringKeys := make(map[string]*Summary)
	err := json.Unmarshal(data, &stringKeys)
	if err != nil {
		return err
	}

	for m, s := range stringKeys {
		month, err := StringToMonth(m)
		if err != nil {
			return err
		}
		ms[month] = s
	}

	return nil
}

func StringToMonth(s string) (time.Month, error) {
	d, err := time.Parse("January", s)
	if err != nil {
		return time.January, err
	}

	return d.Month(), nil
}

func (ms MonthlySummaries) Save(prefix string) error {
	name := fmt.Sprintf("%s.json", prefix)

	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	err = e.Encode(ms)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Not needed now?
func LoadSummaries(prefix string) (MonthlySummaries, error) {
	result := make(MonthlySummaries)

	name := fmt.Sprintf("%s.json", prefix)

	f, err := os.Open(name)
	if err != nil {
		return result, err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
