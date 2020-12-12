package main

import (
	"github.com/Lorenc326/tooth-gnome/locales"
	"log"
	"time"

	"github.com/Lorenc326/tooth-gnome/messages"
	"github.com/Lorenc326/tooth-gnome/orm"

	cron "github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

func handleFatal() {
	if err := recover(); err != nil {
		log.Fatal(err)
	}
}

var watcher = cron.New()

func main() {
	defer handleFatal()

	if err := locales.PreloadLocales("locales/assets"); err != nil {
		log.Fatal(err)
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  config.botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	db := orm.ConnectDB(config.postgresUrl)
	defer db.Close()

	approveMark, approveBtn := messages.BuildApprovalMarkup()
	watcher.AddFunc("* * * * *", messages.GetReminderWatcher(db, bot, approveMark))
	watcher.Start()
	defer watcher.Stop()

	bot.Handle("/start", messages.GetStartHandler(db, bot))
	bot.Handle("/time", messages.GetTimeHandler(db, bot))
	bot.Handle("/progress", messages.GetProgressHandler(db, bot))
	bot.Handle("/language", messages.GetLanguageHandler(db, bot))
	bot.Handle(approveBtn, messages.GetApprovalHandler(db, bot))
	bot.Start()
}
