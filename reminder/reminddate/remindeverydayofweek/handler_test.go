package remindeverydayofweek_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindeverydayofweek"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindEveryDayOfWeek(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindeverydayofweek.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		type TestCase struct {
			Text                       string
			ExpectedRepeatableDateTime *reminder.RepeatableDateTime
		}

		testCases := map[string]TestCase{
			"without hours and minutes": {
				Text: "/remind me every monday buy milk",
				ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
					DayOfWeek: "1",
					Month:     "*",
					Hour:      "9",
					Minute:    "0",
				},
			},
			"with hours and minutes": {
				Text: "/remind me every monday at 23:34 buy milk",
				ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
					DayOfWeek: "1",
					Month:     "*",
					Hour:      "23",
					Minute:    "34",
				},
			},
			"with hours and minutes dot separator": {
				Text: "/remind me every monday at 23.34 buy milk",
				ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
					DayOfWeek: "1",
					Month:     "*",
					Hour:      "23",
					Minute:    "34",
				},
			},
		}

		for name := range testCases {
			t.Run(name, func(t *testing.T) {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				bot := fakeBot.NewBot()
				c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
				mockReminderService := mocks.NewMockServicer(mockCtrl)
				mockReminderService.
					EXPECT().
					AddRepeatableReminderOnDateTime(
						1,
						testCases[name].Text,
						testCases[name].ExpectedRepeatableDateTime,
						"buy milk").
					Return(time.Now(), nil)

				err := remindeverydayofweek.HandleRemindEveryDayOfWeek(mockReminderService)(c)
				require.NoError(t, err)
				require.Len(t, bot.OutboundSendMessages, 1)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		text := "/remind me every monday buy milk"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddRepeatableReminderOnDateTime(
				1,
				text,
				&reminder.RepeatableDateTime{
					DayOfWeek: "1",
					Month:     "*",
					Hour:      "9",
					Minute:    "0",
				},
				"buy milk").
			Return(time.Now(), errors.New("error"))

		err := remindeverydayofweek.HandleRemindEveryDayOfWeek(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
