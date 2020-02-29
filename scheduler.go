package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Scheduler struct {
	db *gorm.DB
}

func (s *Scheduler) Schedule(user *User) {
	if user.StartsAt == nil {
		return
	}
	if user.EveryHours == nil {
		return
	}
	s.db.Where("user_id = ?", user.ID).Unscoped().Delete(&Schedule{})
	start := *user.StartsAt + time.Minute*15
	till := *user.StartsAt + time.Hour*8

	for {
		schedule := Schedule{
			UserID: user.ID,
			At:     start,
		}
		s.db.Save(&schedule)

		start += time.Hour * time.Duration(*user.EveryHours)
		if start > till {
			break
		}
	}

	user.LastWorkout = nil
	user.LastScheduleID = nil
	s.db.Save(user)
}
