package bot

import tb "gopkg.in/tucnak/telebot.v2"

type Context struct {
	Chat     *tb.Chat
	Text     string
	Callback *tb.Callback
	OwnerID  int
	params   map[string]string
}

func (c *Context) Param(key string) string {
	param := c.params[key]

	return param
}
