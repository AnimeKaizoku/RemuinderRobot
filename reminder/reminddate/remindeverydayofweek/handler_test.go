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

type TestCase struct {
	Text                       string
	ExpectedRepeatableDateTime *reminder.RepeatableDateTime
}

func TestHandleRemindEveryDayOfWeek_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindeverydayofweek.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}
	testCases := newTestHandleRemindEveryDayOfWeekTestCases()

	for name := range testCases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			bot := fakeBot.NewTBWrapBot()
			c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
			mockReminderService := mocks.NewMockServicer(mockCtrl)
			mockReminderService.
				EXPECT().
				AddRepeatableReminderOnDateTime(
					1,
					testCases[name].Text,
					testCases[name].ExpectedRepeatableDateTime,
					"update weekly report").
				Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

			err := remindeverydayofweek.HandleRemindEveryDayOfWeek(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func newTestHandleRemindEveryDayOfWeekTestCases() map[string]TestCase {
	return map[string]TestCase{
		"without hours and minutes": {
			Text: "/remind me every monday update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "9",
				Minute:    "0",
			},
		},
		"with hours and minutes": {
			Text: "/remind me every monday at 23:34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "23",
				Minute:    "34",
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me every monday at 23.34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "23",
				Minute:    "34",
			},
		},
		"with only hour": {
			Text: "/remind me every monday at 23 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "23",
				Minute:    "0",
			},
		},
		"with only hour pm": {
			Text: "/remind me every monday at 8pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "0",
			},
		},
		"with hour minute": {
			Text: "/remind me every monday at 8:30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "30",
			},
		},
		"with hour minute dot separator": {
			Text: "/remind me every monday at 8.30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "30",
			},
		},

		"with time of day and without hours and minutes": {
			Text: "/remind me every monday morning update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "9",
				Minute:    "0",
			},
		},
		"with time of day and hours and minutes": {
			Text: "/remind me every monday evening at 23:34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "23",
				Minute:    "34",
			},
		},
		"with time of day and hours and minutes dot separator": {
			Text: "/remind me every monday night at 23.34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "23",
				Minute:    "34",
			},
		},
		"with time of day and only hour": {
			Text: "/remind me every monday evening at 23 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "23",
				Minute:    "0",
			},
		},
		"with time of day and only hour pm": {
			Text: "/remind me every monday night at 8pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "0",
			},
		},
		"with time of day and hour minute": {
			Text: "/remind me every monday evening at 8:30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "30",
			},
		},
		"with time of day and hour minute dot separator": {
			Text: "/remind me every monday night at 8.30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "30",
			},
		},
		"with time of day only": {
			Text: "/remind me every monday evening update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1",
				Month:     "*",
				Hour:      "20",
				Minute:    "0",
			},
		},
	}
}

func TestHandleRemindEveryDayOfWeek_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindeverydayofweek.HandlePattern)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me every monday update weekly report"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	bot := fakeBot.NewTBWrapBot()
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
			"update weekly report").
		Return(reminder.NextScheduleChatTime{}, errors.New("error"))

	err = remindeverydayofweek.HandleRemindEveryDayOfWeek(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
