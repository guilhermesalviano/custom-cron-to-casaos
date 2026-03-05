package utils

import (
	"github.com/go-co-op/gocron"
)

func ScheduleOnDay(s *gocron.Scheduler, day string) *gocron.Scheduler {
	switch day {
	case "Monday":    return s.Every(1).Monday()
	case "Tuesday":   return s.Every(1).Tuesday()
	case "Wednesday": return s.Every(1).Wednesday()
	case "Thursday":  return s.Every(1).Thursday()
	case "Friday":    return s.Every(1).Friday()
	case "Saturday":  return s.Every(1).Saturday()
	default:          return s.Every(1).Sunday()
	}
}