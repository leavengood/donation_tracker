package main

import (
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

func ParseNvpName(name string) (string, int) {
	if strings.HasPrefix(name, "L_") {
		var i int
		for i = 2; i < len(name); i++ {
			if unicode.IsNumber(rune(name[i])) {
				break
			}
		}
		field := name[2:i]
		number, _ := strconv.Atoi(name[i:])

		return field, number
	} else {
		return "", 0
	}
}

type NameValues map[string]string
type NvpResult struct {
	Ack string
	List map[int]NameValues
}

func (n *NvpResult) Successful() bool {
	return strings.Contains(n.Ack, "Success")
}

func ParseNvpData(data string) *NvpResult {
	result := new(NvpResult)
	result.List = make(map[int]NameValues)

	fields := strings.Split(string(data), "&")

	for _, field := range(fields) {
		nv := strings.Split(field, "=")
		name := nv[0]
		value, _ := url.QueryUnescape(nv[1])

		if name == "ACK" {
			result.Ack = value
		} else {
			field, num := ParseNvpName(name)
			if _, found := result.List[num]; !found {
				result.List[num] = make(NameValues)
			}
			result.List[num][field] = value
		}
	}

	return result
}
