package settimezone

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
)

type Message struct {
	TimeZone string `regexpGroup:"timezone"`
}

// nolint:lll
const HandlePattern = `\/settimezone (?P<timezone>.*)`

func HandleSetTimezone(service Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		err := service.SetTimeZone(int(c.ChatID()), message.TimeZone)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Timezone has been updated to: %s", message.TimeZone))
		return err
	}
}
