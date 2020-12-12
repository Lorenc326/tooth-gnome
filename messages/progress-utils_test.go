package messages

import (
	"fmt"
	"testing"
	"time"

	"github.com/Lorenc326/tooth-gnome/orm"
)

func TestBuildStatisticsMessage(t *testing.T) {
	cases := [][]int16{
		{0, 2},
		{20, 40},
		{40, 40},
		{50, 40},
	}
	for _, args := range cases {
		t.Run(fmt.Sprintf("progress %d, maxProgress %d", args[0], args[1]), func(t *testing.T) {
			if len(buildStatisticsMessage(args[0], args[1])) < 1 {
				t.Errorf("Empty string for %v args", args)
			}
		})
	}
}

type userCase struct {
	user     orm.User
	expected int16
}

func TestGetSkippedProgress(t *testing.T) {
	// stub local variable-wrapper around time.Now()
	mockedNow := "2020-10-15T15:00:00+02:00"
	timeNow = func() time.Time {
		t, _ := time.Parse(time.RFC3339, mockedNow)
		return t
	}

	cases := []userCase{
		// empty input
		userCase{orm.User{LastTrained: "", MorningTime: "10:00+10", EveningTime: "19:00+08"}, 0},
		userCase{orm.User{LastTrained: "2020-11-28T00:00:00+02:00", MorningTime: "", EveningTime: "19:00+08"}, 0},
		userCase{orm.User{LastTrained: "2020-11-28T00:00:00+02:00", MorningTime: "10:00+10", EveningTime: ""}, 0},

		// zero skipped
		userCase{orm.User{LastTrained: "2020-10-15T12:00:00+02:00", MorningTime: "10:00+02", EveningTime: "18:00+02"}, 0},
		userCase{orm.User{LastTrained: "2020-10-15T08:00:00+02:00", MorningTime: "10:00+02", EveningTime: "18:00+02"}, 0},
		userCase{orm.User{LastTrained: "2020-10-15T08:22:00+02:00", MorningTime: "14:10+02", EveningTime: "02:00+02"}, 0},
		userCase{orm.User{LastTrained: "2020-10-15T08:10:00+02:00", MorningTime: "16:22+02", EveningTime: "03:30+02"}, 0},
		userCase{orm.User{LastTrained: "2020-10-14T20:40:00+02:00", MorningTime: "10:10+02", EveningTime: "20:40+02"}, 0},
		userCase{orm.User{LastTrained: "2020-10-15T04:10:00+02:00", MorningTime: "04:10+02", EveningTime: "14:00+02"}, 0},

		// received > 0
		userCase{orm.User{LastTrained: "2020-10-14T19:59:59+02:00", MorningTime: "10:00+02", EveningTime: "20:00+02"}, 2},
		userCase{orm.User{LastTrained: "2020-10-15T03:59:59+02:00", MorningTime: "04:00+02", EveningTime: "14:00+02"}, 2},
		userCase{orm.User{LastTrained: "2020-10-14T09:00:00+02:00", MorningTime: "10:00+02", EveningTime: "20:00+02"}, 2},
		userCase{orm.User{LastTrained: "2020-10-13T19:00:00+02:00", MorningTime: "10:00+02", EveningTime: "20:00+02"}, 6},
		userCase{orm.User{LastTrained: "2020-10-12T19:00:00+02:00", MorningTime: "10:00+02", EveningTime: "20:00+02"}, 10},
	}
	for _, c := range cases {
		if progress := getSkippedProgress(&c.user); progress != c.expected {
			t.Errorf("Expected %d, received %d, input %v", c.expected, progress, c.user)
		}
	}
}
