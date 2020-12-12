package messages

import (
	"github.com/Lorenc326/tooth-gnome/orm"
	"github.com/go-pg/pg/v10"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func BuildApprovalMarkup() (*tb.ReplyMarkup, *tb.Btn) {
	mark := tb.ReplyMarkup{}
	btn := mark.Data("Done", "approve")
	mark.Inline(
		mark.Row(btn),
	)
	return &mark, &btn
}

const maxProgress = 42

func GetApprovalHandler(db *pg.DB, bot *tb.Bot) func(_ *tb.Callback) {
	return func(c *tb.Callback) {
		messageSent := time.Unix(c.Message.Unixtime, 0)
		approvalDeadline := messageSent.Add(3 * time.Hour)
		if time.Now().After(approvalDeadline) {
			bot.Respond(c, &tb.CallbackResponse{
				Text: "Oh no! Its too late!",
			})
			bot.Edit(c.Message, "⌛ Oh! Reminder is already expired! ⌛\nYou should approve it within 3hours, next time")
			return
		}

		user := orm.User{ID: c.Sender.ID}
		user.GetTraining(db)

		reduceProgress := getSkippedProgress(&user)
		if reduceProgress > 0 {
			if user.Progress-reduceProgress <= 0 {
				user.Progress = 0
			} else {
				user.Progress = user.Progress - reduceProgress
			}
		}

		if user.Progress < maxProgress {
			user.Progress += 1
		}
		user.LastTrained = time.Now().Format(time.RFC3339)
		user.Train(db)

		bot.Respond(c, &tb.CallbackResponse{
			Text: "Nice!",
		})
		bot.Edit(c.Message, "👏 You're doing great! 👏\nCheck out your /progress")
	}
}
