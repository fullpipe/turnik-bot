package main

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
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

	s.scheduleAtTimeWithRandom(*user.StartsAt+time.Minute*10, user)
	start := *user.StartsAt + time.Hour*time.Duration(*user.EveryHours)
	till := *user.StartsAt + time.Hour*9

	for {
		s.scheduleAtTimeWithRandom(start, user)

		start += time.Hour * time.Duration(*user.EveryHours)
		if start > till {
			break
		}
	}

	user.LastWorkout = nil
	user.LastScheduleID = nil
	s.db.Save(user)
}

func (s *Scheduler) ResetUserSetting(user *User) {
	user.StartsAt = nil
	user.EveryHours = nil
	user.LastWorkout = nil
	user.LastScheduleID = nil
	s.db.Save(user)
	s.db.Where("user_id = ?", user.ID).Unscoped().Delete(&Schedule{})
}

func (s *Scheduler) scheduleAtTimeWithRandom(start time.Duration, user *User) {
	rand.Seed(time.Now().UnixNano())
	schedule := Schedule{
		UserID: user.ID,
		At:     start + time.Minute*time.Duration(rand.Intn(7)),
	}
	s.db.Save(&schedule)
}
