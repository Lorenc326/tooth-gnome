package messages

import (
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

const introMessage = "Hello there!\nI'll make sure your \U0001F9B7 are washed in time, " +
	"just set your wake up/sleep hours. (it works better then alarm, trust me)"

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
