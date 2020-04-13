package date

import (
	"strings"
)

// nolint:gomnd
func ToNumericMonth(month string) int {
	switch strings.ToLower(month) {
	case "january":
		return 1
	case "february":
		return 2
	case "march":
		return 3
	case "april":
		return 4
	case "may":
		return 5
	case "june":
		return 6
	case "july":
		return 7
	case "august":
		return 8
	case "september":
		return 9
	case "october":
		return 10
	case "november":
		return 11
	case "december":
		return 12
	default:
		return 0
	}
}

// nolint:gomnd
func ToNumericDayOfWeek(day string) int {
	switch strings.ToLower(day) {
	case "sunday":
		return 0
	case "monday":
		return 1
	case "tuesday":
		return 2
	case "wednesday":
		return 3
	case "thursday":
		return 4
	case "friday":
		return 5
	case "saturday":
		return 6
	default:
		return 0
	}
}
