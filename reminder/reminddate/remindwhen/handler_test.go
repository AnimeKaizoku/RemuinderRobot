package remindwhen_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindwhen"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindWhen_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindwhen.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	type TestCase struct {
		Text                 string
		ExpectedWordDateTime reminder.WordDateTime
	}

	testCases := map[string]TestCase{
		"without hours and minutes": {
			Text: "/remind me tonight update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   20,
				Minute: 0,
			},
		},
		"with hours and minutes": {
			Text: "/remind me tonight at 19:45 update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   19,
				Minute: 45,
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me tonight at 19.45 update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   19,
				Minute: 45,
			},
		},
		"with only hour": {
			Text: "/remind me tonight at 19 update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   19,
				Minute: 0,
			},
		},
		"with only hour pm": {
			Text: "/remind me tonight at 8pm update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   20,
				Minute: 0,
			},
		},
		"with hour minute": {
			Text: "/remind me tonight at 8:30pm update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   20,
				Minute: 30,
			},
		},
		"with hour minute dot separator": {
			Text: "/remind me tonight at 8.30pm update weekly report",
			ExpectedWordDateTime: reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   20,
				Minute: 30,
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
				AddReminderOnWordDateTime(
					1,
					testCases[name].Text,
					testCases[name].ExpectedWordDateTime,
					"update weekly report").
				Return(reminddate.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

			err := remindwhen.HandleRemindWhen(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func TestHandleRemindWhen_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindwhen.HandlePattern)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me tonight update weekly report"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	bot := fakeBot.NewTBWrapBot()
	c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
	mockReminderService := mocks.NewMockServicer(mockCtrl)
	mockReminderService.
		EXPECT().
		AddReminderOnWordDateTime(
			1,
			text,
			reminder.WordDateTime{
				When:   reminder.Today,
				Hour:   20,
				Minute: 0,
			},
			"update weekly report").
		Return(reminddate.NextScheduleChatTime{}, errors.New("error"))

	err = remindwhen.HandleRemindWhen(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
