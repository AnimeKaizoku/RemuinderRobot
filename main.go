package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/bot"
	"github.com/enrico5b1b4/telegram-bot/db"
)

// nolint:funlen
func main() {
	dbFile := MustGetEnv("DB_FILE")
	telegramBotToken := MustGetEnv("TELEGRAM_BOT_TOKEN")
	allowedChats := parseAllowedChats(MustGetEnv("ALLOWED_CHATS"))

	database, err := db.SetupDB(dbFile, allowedChats)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	botConfig := tbwrap.Config{
		Token:        telegramBotToken,
		AllowedChats: allowedChats,
	}
	telegramBot, err := tbwrap.NewBot(botConfig)
	if err != nil {
		log.Println(err)
		return
	}

	appBot := bot.New(allowedChats, database, telegramBot)
	appBot.Start()
}

func parseAllowedChats(list string) []int {
	sepList := strings.Split(list, ",")
	intList := make([]int, len(sepList))
	var err error

	for i := range sepList {
		intList[i], err = strconv.Atoi(strings.TrimSpace(sepList[i]))
		if err != nil {
			log.Fatalln(err)
		}
	}

	return intList
}

func MustGetEnv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(fmt.Sprintf("%s must be set", name))
	}

	return value
}
