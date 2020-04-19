package reminddaymonth_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddaymonth"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindDayMonth_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(reminddaymonth.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	type TestCase struct {
		Text             string
		ExpectedDateTime reminder.DateTime
	}

	testCases := map[string]TestCase{
		"without hours and minutes": {
			Text: "/remind me on the 4th of march buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       9,
				Minute:     0,
			},
		},
		"with hours and minutes": {
			Text: "/remind me on the 4th of march at 23:34 buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       23,
				Minute:     34,
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me on the 4th of march at 23.34 buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       23,
				Minute:     34,
			},
		},
		"with only hour": {
			Text: "/remind me on the 4th of march at 23 buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       23,
				Minute:     0,
			},
		},
		"with only hour pm": {
			Text: "/remind me on the 4th of march at 8pm buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       20,
				Minute:     0,
			},
		},
		"with hour minute pm": {
			Text: "/remind me on the 4th of march at 8:30pm buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       20,
				Minute:     30,
			},
		},
		"with hour minute pm dot separator": {
			Text: "/remind me on the 4th of march at 8.30pm buy milk",
			ExpectedDateTime: reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       20,
				Minute:     30,
			},
		},
	}

	for name := range testCases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			bot := fakeBot.NewTBWrapBot()
			c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
			mockReminderService := mocks.NewMockServicer(mockCtrl)
			mockReminderService.
				EXPECT().
				AddReminderOnDateTime(
					1,
					testCases[name].Text,
					testCases[name].ExpectedDateTime,
					"buy milk").
				Return(time.Now(), nil)

			err := reminddaymonth.HandleRemindDayMonth(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func TestHandleRemindDayMonth_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(reminddaymonth.HandlePattern)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me on the 4th of march buy milk"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	bot := fakeBot.NewTBWrapBot()
	c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
	mockReminderService := mocks.NewMockServicer(mockCtrl)
	mockReminderService.
		EXPECT().
		AddReminderOnDateTime(
			1,
			text,
			reminder.DateTime{
				DayOfMonth: 4,
				Month:      3,
				Hour:       9,
				Minute:     0,
			},
			"buy milk").
		Return(time.Now(), errors.New("error"))

	err = reminddaymonth.HandleRemindDayMonth(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
