package store

import (
	"fmt"
	"time"
)

type Weekday int

const (
	Monday    = Weekday(time.Monday)
	Tuesday   = Weekday(time.Tuesday)
	Wednesday = Weekday(time.Wednesday)
	Thursday  = Weekday(time.Thursday)
	Friday    = Weekday(time.Friday)
	Saturday  = Weekday(time.Saturday)
	Sunday    = Weekday(time.Sunday)
)

var (
	weekdayString = map[string]Weekday{
		Monday.String():    Monday,
		Tuesday.String():   Tuesday,
		Wednesday.String(): Wednesday,
		Thursday.String():  Thursday,
		Friday.String():    Friday,
		Saturday.String():  Saturday,
		Sunday.String():    Sunday,
	}
)

func (w Weekday) String() string {
	return time.Weekday(w).String()
}

func WeekdayFromString(s string) (Weekday, error) {
	w, ok := weekdayString[s]
	if !ok {
		return Weekday(0), fmt.Errorf("invalid weekday string: %s", s)
	}

	return w, nil
}

type DaySchedule struct {
	Weekday   Weekday
	StartTime time.Time
	EndTime   time.Time
}

type Schedule struct {
	Days []DaySchedule
}

func (s Schedule) IsActive(t time.Time) bool {
	currentDay := Weekday(t.Weekday())

	for _, v := range s.Days {
		if v.Weekday == currentDay {
			return isTimeBetween(t, v.StartTime, v.EndTime)
		}
	}

	return false
}

func (s Schedule) NextWorkingTime(t time.Time) time.Time {
	if s.IsActive(t) {
		return t
	}

	for t = t.AddDate(0, 0, 1); ; t = t.AddDate(0, 0, 1) {
		weekDay := Weekday(t.Weekday())

		for _, v := range s.Days {
			if v.Weekday == weekDay {
				return time.Date(
					t.Year(), t.Month(), t.Day(),
					v.StartTime.Hour(), v.StartTime.Minute(), v.StartTime.Second(),
					0, v.StartTime.Location(),
				)
			}
		}
	}
}

func isTimeBetween(t time.Time, start time.Time, end time.Time) bool {
	var (
		tSecs     = timeUTCSecs(t)
		startSecs = timeUTCSecs(start)
		endSects  = timeUTCSecs(end)
	)

	if tSecs >= startSecs && tSecs <= endSects {
		return true
	}

	return false
}

func timeUTCSecs(t time.Time) int {
	h, m, s := t.UTC().Clock()
	return h*3600 + m*60 + s
}
