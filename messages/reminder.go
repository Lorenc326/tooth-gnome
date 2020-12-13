package messages

import (
	"github.com/Lorenc326/tooth-gnome/locales"
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

const bunchSize = 100

func GetReminderWatcher(db *pg.DB, bot *tb.Bot, approvalMarkup *tb.ReplyMarkup) func() {
	return func() {
		now := time.Now().Round(time.Minute).Format(reminderTimeFormat)
		log.Println("Reminder watcher", now)
		userModel := new(orm.User)

		firstLoop := true
		processed := 0
		var users []orm.User

		for firstLoop || len(users) >= bunchSize {
			firstLoop = false
			if err := userModel.GetUsersToRemind(db, &users, now, processed, bunchSize); err != nil {
				log.Println(err)
				return
			}
			log.Println("Reminder users", users)
			for _, u := range users {
				bot.Send(u, locales.Translate(u.Lng, "timeToBrush"), approvalMarkup)
				processed += 1
			}
		}
	}
}
