package messages

import (
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"time"
)

const reminderMessage = "Hello! Reminding you to brush your teeth."
const bunchSize = 100

func GetReminderWatcher(db *pg.DB, bot *tb.Bot) func() {
	return func() {
		now := time.Now().Round(time.Minute).Format(reminderTimeFormat)
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
			for _, u := range users {
				bot.Send(u, reminderMessage)
				processed += 1
			}
		}
	}
}
