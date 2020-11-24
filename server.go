package main

import (
	"log"
	"time"

	dbUtils "github.com/Lorenc326/tooth-gnome/db"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	bot, err := tb.NewBot(tb.Settings{
		Token: config.botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		bot.Send(m.Sender, "Hello there!\nI'll make sure your \U0001F9B7 are washed in time, just set your wake up/sleep hours. (it works better then alarm, trust me)")
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Println(m.Text)
	})

	_, err = dbUtils.ConnectDB(config.postgresUrl)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
