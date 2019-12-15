package bot

import (
	"fmt"
	"regexp"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type HandlerFunc func(c *Context) error

type Route struct {
	Path    *regexp.Regexp
	Handler HandlerFunc
}

type Bot struct {
	tBot   *tb.Bot
	routes map[*regexp.Regexp]*Route
}

type Config struct {
	Token         string
	AllowedUsers  []int
	AllowedGroups []int
}

func NewBot(cfg Config) (*Bot, error) {
	poller := NewPollerWithAllowedUserAndGroups(15*time.Second, cfg.AllowedUsers, cfg.AllowedGroups)
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: poller,
	})
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		tBot:   b,
		routes: map[*regexp.Regexp]*Route{},
	}

	// Add default OnText handler
	bot.handle(tb.OnText, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		fmt.Println("OnText")
		bot.HandleOnText(m.Text, m.Chat)
	})

	return bot, nil
}

func (b *Bot) RegisterHandlers() {
}

func (b *Bot) Handle(endpoint interface{}, handler interface{}) {
	b.tBot.Handle(endpoint, handler)
}

func (b *Bot) handle(endpoint interface{}, handler interface{}) {
	b.tBot.Handle(endpoint, handler)
}

func (b *Bot) Respond(callback *tb.Callback, responseOptional ...*tb.CallbackResponse) error {
	return b.tBot.Respond(callback, responseOptional...)
}

func (b *Bot) Send(to tb.Recipient, what interface{}, options ...interface{}) (*tb.Message, error) {
	mergedOptions := append([]interface{}{&tb.SendOptions{ParseMode: tb.ModeMarkdown}}, options...)
	fmt.Printf("%#v\n", mergedOptions)
	fmt.Println("len: ", len(mergedOptions))
	return b.tBot.Send(to, what, mergedOptions...)
}

func (b *Bot) AddRegExp(path string, handler HandlerFunc) {
	compiledRegExp := regexp.MustCompile(path)

	b.routes[compiledRegExp] = &Route{Path: compiledRegExp, Handler: handler}
}

func (b *Bot) AddMultiRegExp(paths []string, handler HandlerFunc) {
	for i := range paths {
		compiledRegExp := regexp.MustCompile(paths[i])

		b.routes[compiledRegExp] = &Route{Path: compiledRegExp, Handler: handler}
	}
}

func (b *Bot) Add(path string, handler HandlerFunc) {
	b.handle(path, func(m *tb.Message) {
		fmt.Println("in add() for path ", path)
		c := &Context{Chat: m.Chat, Text: m.Text, OwnerID: int(m.Chat.ID)}
		handler(c)
		return
	})
}

// TODO
func (b *Bot) AddButton(path *tb.InlineButton, handler HandlerFunc) {
	b.handle(path, func(callback *tb.Callback) {
		fmt.Println("in AddButton() for button ", path)
		c := &Context{
			Chat:     callback.Message.Chat,
			Text:     callback.Message.Text,
			Callback: callback,
			OwnerID:  int(callback.Message.Chat.ID),
		}
		handler(c)
		return
	})
}

func (b *Bot) HandleOnText(text string, chat *tb.Chat) {
	for regExpKey := range b.routes {
		matches := regExpKey.FindStringSubmatch(text)
		names := regExpKey.SubexpNames()

		if len(matches) > 0 {
			fmt.Printf("found match with regex:  %v\n", regExpKey)
			params := mapSubexpNames(matches, names)
			c := &Context{Chat: chat, Text: text, params: params, OwnerID: int(chat.ID)}
			b.routes[regExpKey].Handler(c)
			return

			//return i, mapSubexpNames(matches, names), nil
			// for k := range mapNames {
			// 	if k != "" {
			// 		fmt.Printf("%s: %s\n", k, mapNames[k])
			// 	}
			// }
			// break
		}
	}
}

func (b *Bot) Start() {
	// And this one — just under the message itself.
	// Pressing it will cause the client to send
	// the bot a callback.
	//
	// Make sure Unique stays unique as it has to be
	// for callback routing to work.
	confirmRemindDeleteBtn := tb.InlineButton{
		Unique: "confirmRemindDeleteBtn",
		Text:   "Yes",
		Data:   "1",
	}
	cancelRemindDeleteBtn := tb.InlineButton{
		Unique: "cancelRemindDeleteBtn",
		Text:   "No",
		Data:   "1",
	}
	remindDeleteInlineKeys := [][]tb.InlineButton{
		{confirmRemindDeleteBtn},
		{cancelRemindDeleteBtn},
	}

	b.tBot.Handle(&confirmRemindDeleteBtn, func(c *tb.Callback) {
		// TODO keep global map of last item to delete
		fmt.Println("c.Data")
		fmt.Println(c.Data)
		b.tBot.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		b.Send(c.Sender, "reminder deleted")

		// on inline button pressed (callback!)
		//b.tBot.Send(
		//	c.Sender,
		//	`You're edit item send me new name for it.`,
		//)
		//
		//// always respond!
		//b.tBot.Respond(c, &tb.CallbackResponse{
		//	CallbackID: "123",
		//	Text:       "asd",
		//	ShowAlert:  true,
		//	URL:        "/remindlist",
		//})

	})

	// Command: /start <PAYLOAD>
	b.Handle("/enrico1", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		// TODO keep global map of last item to delete
		//b.Send()

		b.tBot.Send(m.Sender, `(inline URL)[/remindlist]`, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
		remindDeleteInlineKeys[0][0].Data = "1"
		b.tBot.Send(m.Sender, "Do you really want to delete reminder number X", &tb.ReplyMarkup{
			InlineKeyboard:      remindDeleteInlineKeys,
			ReplyKeyboardRemove: true,
		})
	})

	b.Handle("/enrico2", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		// TODO keep global map of last item to delete
		//b.Send()

		b.tBot.Send(m.Sender, `(inline URL)[/remindlist]`, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
		remindDeleteInlineKeys[0][0].Data = "2"
		b.tBot.Send(m.Sender, "Do you really want to delete reminder number X", &tb.ReplyMarkup{
			InlineKeyboard:      remindDeleteInlineKeys,
			ReplyKeyboardRemove: true,
		})
	})

	b.tBot.Start()
}

func mapSubexpNames(m, n []string) map[string]string {
	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i := range n {
		r[n[i]] = m[i]
	}
	return r
}
