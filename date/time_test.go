package date_test

import (
	"fmt"
	"testing"

	"github.com/enrico5b1b4/telegram-bot/date"
	"github.com/stretchr/testify/assert"
)

func TestConvertTo24H(t *testing.T) {
	type TestCase struct {
		Hour           int
		Minute         int
		AMPM           string
		ExpectedHour   int
		ExpectedMinute int
	}

	testCases := []TestCase{
		{
			Hour:           1,
			Minute:         30,
			AMPM:           "pm",
			ExpectedHour:   13,
			ExpectedMinute: 30,
		},
		{
			Hour:           12,
			Minute:         0,
			AMPM:           "pm",
			ExpectedHour:   12,
			ExpectedMinute: 0,
		},
		{
			Hour:           12,
			Minute:         30,
			AMPM:           "am",
			ExpectedHour:   0,
			ExpectedMinute: 30,
		},
		{
			Hour:           1,
			Minute:         0,
			AMPM:           "am",
			ExpectedHour:   1,
			ExpectedMinute: 0,
		},
		{
			Hour:           13,
			Minute:         0,
			AMPM:           "",
			ExpectedHour:   13,
			ExpectedMinute: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%d-%s", tc.Hour, tc.Minute, tc.AMPM), func(t *testing.T) {
			hour, minute := date.ConvertTo24H(tc.Hour, tc.Minute, tc.AMPM)

			assert.Equal(t, tc.ExpectedHour, hour)
			assert.Equal(t, tc.ExpectedMinute, minute)
		})
	}
}
