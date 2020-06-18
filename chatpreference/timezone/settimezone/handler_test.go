package settimezone_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/chatpreference/timezone/settimezone"
	"github.com/enrico5b1b4/telegram-bot/chatpreference/timezone/settimezone/mocks"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleSetTimezone(t *testing.T) {
	handlerPattern, err := regexp.Compile(settimezone.HandlePattern)
	require.NoError(t, err)
	text := "/settimezone Europe/London"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockService := mocks.NewMockServicer(mockCtrl)
		mockService.
			EXPECT().
			SetTimeZone(1, "Europe/London").
			Return(nil)

		err := settimezone.HandleSetTimezone(mockService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockService := mocks.NewMockServicer(mockCtrl)
		mockService.
			EXPECT().
			SetTimeZone(1, "Europe/London").
			Return(errors.New("error"))

		err := settimezone.HandleSetTimezone(mockService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
