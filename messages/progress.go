package messages

import (
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

const progressMessageHeader = "📈 Your habit progress 📈"

func buildStatistics(progress int16, maxProgress int16) string {
	message := make([]string, 200)
	for i := int16(0); i < maxProgress; i += 2 {
		day := i/2 + 1
		if day == 1 {
			message = append(message, "\n\nWeek 1:")
		} else if day == 8 {
			message = append(message, "\n\nWeek 2:")
		} else if day == 15 {
			message = append(message, "\n\nWeek 3:")
		}

		// two points == 1 day
		// 1 success day - green square 🟩
		// 0.5 success day - yellow square 🟨
		// 1 unreached day - blue square 🟦
		morningDone := progress >= i+1
		eveningDone := progress >= i+2
		if morningDone && eveningDone {
			message = append(message, " \U0001F7E9")
		} else if morningDone || eveningDone {
			message = append(message, " \U0001F7E8")
		} else {
			message = append(message, " \U0001F7E6")
		}
	}
	return strings.Join(message, "")
}

func GetProgressHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		user := &orm.User{ID: m.Sender.ID}
		user.GetTraining(db)

		message := progressMessageHeader + buildStatistics(user.Progress, maxProgress)
		bot.Send(m.Sender, message)
	}
}
