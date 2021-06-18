package humanize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func TimeSince(a time.Time) string {
	return formatDiff(diff(a, time.Now()))
}

func StringToSeconds(s string) (int, error) {
	if strings.HasSuffix(s, "m") {
		num := strings.TrimSuffix(s, "m")
		integer, err := strconv.Atoi(num)

		return integer * 60, err
	}

	if strings.HasSuffix(s, "s") {
		num := strings.TrimSuffix(s, "s")
		integer, err := strconv.Atoi(num)

		return integer, err
	}

	integer, err := strconv.Atoi(s)

	return integer, err
}

func SecondsToString(s int) string {
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	if s%60 == 0 {
		return fmt.Sprintf("%dm", s/60)
	}

	floored := math.Floor(float64(s) / 60)
	rest := float64(s) - (floored * 60)

	return fmt.Sprintf("%fm %fs", floored, rest)
}

func formatDiff(years, months, days, hours, mins, secs int) string {
	since := ""
	if years > 0 {
		switch years {
		case 1:
			since += fmt.Sprintf("%d year ", years)
		default:
			since += fmt.Sprintf("%d years ", years)
		}
	}
	if months > 0 {
		switch months {
		case 1:
			since += fmt.Sprintf("%d month ", months)
		default:
			since += fmt.Sprintf("%d months ", months)
		}
	}
	if days > 0 {
		switch days {
		case 1:
			since += fmt.Sprintf("%d day ", days)
		default:
			since += fmt.Sprintf("%d days ", days)
		}
	}
	if hours > 0 {
		switch hours {
		case 1:
			since += fmt.Sprintf("%d hour ", hours)
		default:
			since += fmt.Sprintf("%d hours ", hours)
		}
	}
	if mins > 0 && days == 0 && months == 0 && years == 0 {
		switch mins {
		case 1:
			since += fmt.Sprintf("%d min ", mins)
		default:
			since += fmt.Sprintf("%d mins ", mins)
		}
	}
	if secs > 0 && days == 0 && months == 0 && years == 0 && hours == 0 {
		switch secs {
		case 1:
			since += fmt.Sprintf("%d sec ", secs)
		default:
			since += fmt.Sprintf("%d secs ", secs)
		}
	}
	return strings.TrimSpace(since)
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}
	return
}
