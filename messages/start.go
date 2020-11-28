package messages

import (
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

const introMessage = "👋 Hello there!\n\n" +
	"Develop a habit to brush \U0001F9B7 in 21 days! I will send you reminders twice a day, you can also check your progress any time.\n\n" +
	"Type when you would like to receive reminders +timezone.\n" +
	"Example \"/time 09:00 21:20 +02 \" or \"/time 11:00 23:00 -07 \""

func GetStartHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		user := &orm.User{
			ID:        m.Sender.ID,
			Lng:       m.Sender.LanguageCode,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		user.InsertIfNotExist(db)

		bot.Send(m.Sender, introMessage)
	}
}
