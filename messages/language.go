package messages

import (
	"github.com/Lorenc326/tooth-gnome/locales"
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
)

const languageMessage = "🇯🇵 S 🇰🇷 U 🇩🇪 C 🇨🇳 C 🇺🇸 E 🇫🇷 S 🇪🇸 S 🇬🇧"

var supportedLanguages = []string{"uk", "en", "ru"}

func isLanguageSupported(input string) bool {
	for _, lng := range supportedLanguages {
		if lng == input {
			return true
		}
	}
	return false
}

func GetLanguageHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		user := &orm.User{
			ID: m.Sender.ID,
		}
		user.GetTraining(db)

		if !isLanguageSupported(m.Payload) {
			bot.Send(m.Sender, locales.Translate(user.Lng, "invalidInputLanguage"))
			return
		}

		user = &orm.User{
			ID:  m.Sender.ID,
			Lng: m.Payload,
		}
		user.UpdateLng(db)

		bot.Send(m.Sender, languageMessage)
	}
}
