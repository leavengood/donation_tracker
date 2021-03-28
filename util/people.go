package util

import "sort"

type Person struct {
	Name  string
	Email string
	Total CurrencyAmounts
	Count int
}

type People []*Person

func (p People) Len() int      { return len(p) }
func (p People) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

type ByTotal struct{ People }

func (s ByTotal) Less(i, j int) bool {
	totalI := s.People[i].Total.GrandTotal(1)
	totalJ := s.People[j].Total.GrandTotal(1)
	// We want descending order
	return totalI >= totalJ
}

func (p People) Sort() {
	sort.Sort(ByTotal{p})
}
