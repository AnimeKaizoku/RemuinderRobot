package remindin_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindin"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

// nolint:dupl
func TestHandleRemindInPattern1(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindin.HandlePattern1)
	require.NoError(t, err)
	text := "/remind me in 2 minutes update weekly report"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderIn(
				1,
				text,
				reminder.AmountDateTime{
					Days:    0,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminddate.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

		err := remindin.HandleRemindIn(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderIn(
				1,
				text,
				reminder.AmountDateTime{
					Days:    0,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminddate.NextScheduleChatTime{}, errors.New("error"))

		err := remindin.HandleRemindIn(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}

// nolint:dupl
func TestHandleRemindInPattern2(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindin.HandlePattern2)
	require.NoError(t, err)
	text := "/remind me in 2 minutes, 3 days update weekly report"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderIn(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminddate.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

		err := remindin.HandleRemindIn(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderIn(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   0,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminddate.NextScheduleChatTime{}, errors.New("error"))

		err := remindin.HandleRemindIn(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}

// nolint:dupl
func TestHandleRemindInPattern3(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindin.HandlePattern3)
	require.NoError(t, err)
	text := "/remind me in 2 minutes, 1 hour, 3 days update weekly report"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderIn(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   1,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminddate.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

		err := remindin.HandleRemindIn(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderIn(
				1,
				text,
				reminder.AmountDateTime{
					Days:    3,
					Hours:   1,
					Minutes: 2,
				},
				"update weekly report").
			Return(reminddate.NextScheduleChatTime{}, errors.New("error"))

		err := remindin.HandleRemindIn(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
