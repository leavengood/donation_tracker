package util

import "time"

const (
	dateFormat     = "Jan 2, 2006"
	dateTimeFormat = "Jan 2, 2006 3:05 pm UTC"
)

func FormatDate(t time.Time) string {
	return t.Format(dateFormat)
}

func FormatDateTime(t time.Time) string {
	return t.Format(dateTimeFormat)
}
