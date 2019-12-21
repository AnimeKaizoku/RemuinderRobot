package remindhelp_test

import (
	"regexp"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder/remindhelp"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindHelp(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindhelp.HandlePattern)
	require.NoError(t, err)
	text := "/remindhelp"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)

		err := remindhelp.HandleRemindHelp()(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})
}
