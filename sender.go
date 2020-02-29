package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Sender struct {
	bot *tb.Bot
	db  *gorm.DB
}

func (s *Sender) SendMotivations() {
	schedules := []Schedule{}
	now := time.Now()
	s.db.Where("at BETWEEN ? AND ?", FromBod(now)-time.Minute*2, FromBod(now)).Find(&schedules)
	for _, schedule := range schedules {
		var user User
		s.db.Model(&schedule).Related(&user)

		if user.LastWorkout != nil && now.Sub(*user.LastWorkout) < time.Hour*1 {
			continue
		}

		if user.LastWorkout != nil && user.LastScheduleID == &schedule.ID {
			continue
		}
		log.Println(user.Recipient())
		s.bot.Send(&user, "Пора бы подтянуться")
		user.LastWorkout = &now
		user.LastScheduleID = &schedule.ID
		s.db.Save(user)
	}
}
