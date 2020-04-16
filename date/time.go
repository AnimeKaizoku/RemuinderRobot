package date

//nolint:gocritic
func ConvertTo24H(hour, minute int, ampm string) (int, int) {
	// midnight
	if ampm == "am" && hour == 12 {
		return 0, minute
	}

	// noon
	if ampm == "pm" && hour == 12 {
		//nolint:gomnd
		return 12, minute
	}

	if ampm == "pm" {
		return hour + 12, minute
	}

	return hour, minute
}
