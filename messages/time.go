package messages

import (
	"fmt"
	"github.com/Lorenc326/tooth-gnome/locales"
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"regexp"
	"strings"
	"time"
)

const timeMessage = "👍"
const reminderTimeFormat = "15:04Z07"

func debugTimeHandler(m *tb.Message, bot *tb.Bot) {
	if err := recover(); err != nil {
		bot.Send(m.Sender, locales.Translate(m.Sender.LanguageCode, "invalidInputTime"))
	}
}

func GetTimeHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}
		defer debugTimeHandler(m, bot)

		user := &orm.User{ID: m.Sender.ID}
		user.GetTraining(db)

		reg := regexp.MustCompile("[\\s\\,]+")
		// Expected payload in "10:00 21:00 +2" format
		times := reg.Split(strings.TrimSpace(m.Payload), 3)
		if len(times) != 3 {
			bot.Send(m.Sender, locales.Translate(user.Lng, "invalidInputTime"))
			return
		}

		zoneOffset := times[2]
		layout := "15:04 MST -07"
		start, err1 := time.Parse(layout, fmt.Sprintf("%s GMT %s", times[0], zoneOffset))
		end, err2 := time.Parse(layout, fmt.Sprintf("%s GMT %s", times[1], zoneOffset))
		if err1 != nil || err2 != nil || start.After(end) {
			bot.Send(m.Sender, locales.Translate(user.Lng, "invalidInputTime"))
			return
		}

		user = &orm.User{
			ID:          m.Sender.ID,
			MorningTime: start.Format(reminderTimeFormat),
			EveningTime: end.Format(reminderTimeFormat),
		}
		// TODO: do smth in case of failed save
		user.SetReminders(db)

		bot.Send(m.Sender, timeMessage)
	}
}
