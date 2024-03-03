package helper

import "time"

func StringToDate(s string) time.Time {
	date, _ := time.Parse("2006-01-02", s)
	return date
}

func DateToString(d time.Time) string {
	return d.Format("2006-01-02")
}
