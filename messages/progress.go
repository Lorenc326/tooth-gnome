package messages

import (
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
)

func GetProgressHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		user := &orm.User{ID: m.Sender.ID}
		user.GetTraining(db)

		reduceProgress := getSkippedProgress(user)

		message := progressMessageHeader + buildStatisticsMessage(user.Progress-reduceProgress, maxProgress)
		bot.Send(m.Sender, message)
	}
}
