package messages

import (
	"github.com/Lorenc326/tooth-gnome/locales"
	"github.com/Lorenc326/tooth-gnome/orm"
	"math"
	"strings"
	"time"
)

func buildStatisticsMessage(lng string, progress int16, maxProgress int16) string {
	message := make([]string, 200)
	for i := int16(0); i < maxProgress; i += 2 {
		day := i/2 + 1
		if day == 1 {
			message = append(message, "\n\n"+locales.Translate(lng, "week1"))
		} else if day == 8 {
			message = append(message, "\n\n"+locales.Translate(lng, "week2"))
		} else if day == 15 {
			message = append(message, "\n\n"+locales.Translate(lng, "week3"))
		}
		// two points == 1 day
		morningDone := progress >= i+1
		eveningDone := progress >= i+2
		if morningDone && eveningDone {
			// 1 success day - green square 🟩
			message = append(message, " \U0001F7E9")
		} else if morningDone || eveningDone {
			// 0.5 success day - yellow square 🟨
			message = append(message, " \U0001F7E8")
		} else {
			// 1 unreached day - blue square 🟦
			message = append(message, " \U0001F7E6")
		}
	}
	return strings.Join(message, "")
}

// parses hour of reminder format 10:00+02 to time
func parseHourStamp(hourStamp string) (*time.Time, error) {
	stamp, err := time.Parse(reminderTimeFormat, hourStamp)
	if err != nil {
		return &stamp, err
	}
	// add 1 day just to not overlap 0000 year, date anyway doesn't matter
	parsed := stamp.Add(24 * time.Hour).UTC()
	return &parsed, nil
}

// define as a package variable so it can be stubbed in tests
var timeNow = func() time.Time {
	return time.Now()
}

// get amount of points that should be reduced from progress
// happens when user skips reminders
// keep all in UTC to have unified locations
func getSkippedProgress(user *orm.User) int16 {
	// proceed with calc only if user has trained already
	if user.LastTrained == "" || user.EveningTime == "" || user.MorningTime == "" {
		return 0
	}
	lastTrainedBase, _ := time.Parse(time.RFC3339, user.LastTrained)
	lastTrained := lastTrainedBase.UTC()
	now := timeNow().UTC()
	morning, _ := parseHourStamp(user.MorningTime)
	evening, _ := parseHourStamp(user.EveningTime)

	// we don't care about morning or evening, we need lower and upper hour bounds
	// if after utc formatting evening hour is lower - we swap them
	if evening.Hour() < morning.Hour() || (evening.Hour() == morning.Hour() && evening.Minute() <= morning.Minute()) {
		morning, evening = evening, morning
	}

	// daytime vs nighttime cases
	var checkReminder time.Time
	afterMorning := now.Hour() > morning.Hour() || (now.Hour() == morning.Hour() && now.Minute() >= morning.Minute())
	beforeEvening := now.Hour() < evening.Hour() || (now.Hour() == evening.Hour() && now.Minute() < evening.Minute())
	if !afterMorning && beforeEvening {
		checkReminder = time.
			Date(now.Year(), now.Month(), now.Day(), morning.Hour(), morning.Minute(), 0, 0, time.UTC).
			AddDate(0, 0, -1)
	} else if afterMorning && beforeEvening {
		checkReminder = time.
			Date(now.Year(), now.Month(), now.Day(), evening.Hour(), evening.Minute(), 0, 0, time.UTC).
			AddDate(0, 0, -1)
	} else {
		checkReminder = time.Date(now.Year(), now.Month(), now.Day(), morning.Hour(), morning.Minute(), 0, 0, time.UTC)
	}

	// if at least one reminder was skipped
	// we truncate 2 points + 4 points for each next skipped day
	if lastTrained.UTC().Before(checkReminder) {
		sub := checkReminder.Sub(lastTrained.UTC())
		return int16(2 + math.Floor(sub.Hours()/24)*4)
	}
	return 0
}
