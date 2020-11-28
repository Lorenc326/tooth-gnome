package messages

import (
	"fmt"
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"regexp"
	"strings"
	"time"
)

const timeMessage = "👍"

func GetTimeHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		reg := regexp.MustCompile("[\\s\\,]+")
		// Expected payload in "10:00 21:00 +2" format
		times := reg.Split(strings.TrimSpace(m.Payload), 3)
		zoneOffset := times[2]
		layout := "15:04 MST -07"
		start, _ := time.Parse(layout, fmt.Sprintf("%s GMT %s", times[0], zoneOffset))
		end, _ := time.Parse(layout, fmt.Sprintf("%s GMT %s", times[1], zoneOffset))

		user := &orm.User{
			ID:          m.Sender.ID,
			MorningTime: start.String(),
			EveningTime: end.String(),
		}
		user.SetReminders(db)

		bot.Send(m.Sender, timeMessage)
	}
}
