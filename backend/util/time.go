package util

import "time"

func GetUTCDateTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func FormatDateTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
