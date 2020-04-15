package reminddaymonthyear_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddaymonthyear"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindDayMonthYear(t *testing.T) {
	handlerPattern, err := regexp.Compile(reminddaymonthyear.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		type TestCase struct {
			Text             string
			ExpectedDateTime reminder.DateTime
		}

		testCases := map[string]TestCase{
			"without hours and minutes": {
				Text: "/remind me on the 4th of march 2020 buy milk",
				ExpectedDateTime: reminder.DateTime{
					Day:    4,
					Month:  3,
					Year:   2020,
					Hour:   9,
					Minute: 0,
				},
			},
			"with hours and minutes": {
				Text: "/remind me on the 4th of march 2020 at 23:34 buy milk",
				ExpectedDateTime: reminder.DateTime{
					Day:    4,
					Month:  3,
					Year:   2020,
					Hour:   23,
					Minute: 34,
				},
			},
			"with hours and minutes dot separator": {
				Text: "/remind me on the 4th of march 2020 at 23.34 buy milk",
				ExpectedDateTime: reminder.DateTime{
					Day:    4,
					Month:  3,
					Year:   2020,
					Hour:   23,
					Minute: 34,
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
					AddReminderOnDateTime(
						1,
						testCases[name].Text,
						testCases[name].ExpectedDateTime,
						"buy milk").
					Return(time.Now(), nil)

				err := reminddaymonthyear.HandleRemindDayMonthYear(mockReminderService)(c)
				require.NoError(t, err)
				require.Len(t, bot.OutboundSendMessages, 1)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		text := "/remind me on the 4th of march 2020 buy milk"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderOnDateTime(
				1,
				text,
				reminder.DateTime{
					Day:    4,
					Month:  3,
					Year:   2020,
					Hour:   9,
					Minute: 0,
				},
				"buy milk").
			Return(time.Now(), errors.New("error"))

		err := reminddaymonthyear.HandleRemindDayMonthYear(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
