package util

import (
	"fmt"
	"time"
)

func FormatAgeFromNow(from time.Time) string {
	return FormatAgeFrom(from, time.Now())
}

func FormatAgeFrom(from, to time.Time) string {
	delta := from.Sub(to)
	switch {
	case delta >= 0:
		return FormatAge(delta)
	default:
		return FormatAge(-delta)
	}
}

func FormatAge(d time.Duration) string {
	switch {
	case d >= time.Hour*24:
		return fmt.Sprintf("%dd", d/(time.Hour*24))
	case d >= time.Hour:
		return fmt.Sprintf("%dh", d/time.Hour)
	case d >= time.Minute:
		return fmt.Sprintf("%dm", d/time.Minute)
	case d >= time.Second:
		return fmt.Sprintf("%ds", d/time.Second)
	default:
		return fmt.Sprintf("%0.1fs", float64(d)/float64(time.Second))
	}
}
