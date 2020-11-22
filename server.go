package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	fmt.Println("Token:", os.Getenv("bot_token"))

	bot, err := tb.NewBot(tb.Settings{
		Token: "asd",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Println(m.Text)
	})

	bot.Start()
}
