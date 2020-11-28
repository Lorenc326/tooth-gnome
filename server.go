package main

import (
	"log"
	"time"

	"github.com/Lorenc326/tooth-gnome/orm"

	tb "gopkg.in/tucnak/telebot.v2"
)

func handleFatal() {
	if err := recover(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer handleFatal()

	bot, err := tb.NewBot(tb.Settings{
		Token:  config.botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	db := orm.ConnectDB(config.postgresUrl)
	defer db.Close()

	bot.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}

		user := &orm.User{ID: m.Sender.ID, Lng: m.Sender.LanguageCode, CreatedAt: time.Now().Format(time.RFC3339)}
		db.Model(user).SelectOrInsert()

		bot.Send(m.Sender, "Hello there!\nI'll make sure your \U0001F9B7 are washed in time, just set your wake up/sleep hours. (it works better then alarm, trust me)")
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		log.Println(m.Text, m.Sender, m.Sender.Username)
	})

	bot.Start()
}
