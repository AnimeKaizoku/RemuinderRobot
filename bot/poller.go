package bot

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func NewPollerWithAllowedUserAndGroups(pollTimout time.Duration, allowedUsers []int, allowedGroups []int) *tb.MiddlewarePoller {
	poller := &tb.LongPoller{Timeout: pollTimout}
	return tb.NewMiddlewarePoller(poller, func(upd *tb.Update) bool {
		// TODO check how groups work
		fmt.Printf("%#v\n", upd)

		allowedUsersAndGroups := append(allowedUsers, allowedGroups...)

		// TODO support edited messages - upd.EditedMessage
		if upd.Message != nil {
			fmt.Printf("%#v\n", upd.Message)
			fmt.Println(upd.Message.Text)
			fmt.Println("upd.Message.Sender")
			fmt.Println(upd.Message.Sender.ID)
			fmt.Println(upd.Message.Sender.FirstName)
			fmt.Println("upd.Message.Chat")
			fmt.Println(upd.Message.Chat.ID)
			fmt.Println(upd.Message.Chat.Title)

			if !isInList(int(upd.Message.Chat.ID), allowedUsersAndGroups) {
				return false
			}
			return true
		}

		if upd.Callback != nil {
			fmt.Printf("%#v\n", upd.Callback)
			fmt.Println(upd.Callback.ID)
			fmt.Println("upd.Callback.Data")
			fmt.Println(upd.Callback.Data)
			fmt.Println("upd.Callback.Message.Text")
			fmt.Println(upd.Callback.Message.Text)
			fmt.Println("upd.Callback.Message.Sender.ID")
			fmt.Println(upd.Callback.Message.Sender.ID)
			fmt.Println("upd.Callback.Message.Chat.ID")
			fmt.Println(upd.Callback.Message.Chat.ID)
			fmt.Println("upd.Callback.Sender.ID")
			fmt.Println(upd.Callback.Sender.ID)

			if !isInList(int(upd.Callback.Message.Chat.ID), allowedUsersAndGroups) {
				return false
			}
			return true
		}

		return false
	})
}

func isInList(ID int, list []int) bool {
	for i := range list {
		if list[i] == ID {
			return true
		}
	}
	return false
}
