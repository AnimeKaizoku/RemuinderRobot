package reminddayofweek_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/reminddayofweek"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindDayOfWeek_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(reminddayofweek.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	type TestCase struct {
		Text             string
		ExpectedDateTime reminder.DateTime
	}

	testCases := map[string]TestCase{
		"without hours and minutes": {
			Text: "/remind me on tuesday update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      9,
				Minute:    0,
			},
		},
		"with hours and minutes": {
			Text: "/remind me on tuesday at 23:34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    34,
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me on tuesday at 23.34 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    34,
			},
		},
		"with only hour": {
			Text: "/remind me on tuesday at 23 update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      23,
				Minute:    0,
			},
		},
		"with only hour pm": {
			Text: "/remind me on tuesday at 8pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    0,
			},
		},
		"with hour minute pm": {
			Text: "/remind me on tuesday at 8:30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    30,
			},
		},
		"with hour minute pm dot separator": {
			Text: "/remind me on tuesday at 8.30pm update weekly report",
			ExpectedDateTime: reminder.DateTime{
				DayOfWeek: "2",
				Hour:      20,
				Minute:    30,
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
					"update weekly report").
				Return(reminddate.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

			err := reminddayofweek.HandleRemindDayOfWeek(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func TestHandleRemindDayOfWeek_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(reminddayofweek.HandlePattern)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me on tuesday update weekly report"
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
				DayOfWeek: "2",
				Hour:      9,
				Minute:    0,
			},
			"update weekly report").
		Return(reminddate.NextScheduleChatTime{}, errors.New("error"))

	err = reminddayofweek.HandleRemindDayOfWeek(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
