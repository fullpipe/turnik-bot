package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Motivator struct {
	bot *tb.Bot
	db  *gorm.DB
}

type Motivation struct {
	Text string
	URL  string
}

var motivations = []Motivation{
	Motivation{
		Text: "Мэнни хочет чтобы ты подтянулся",
		URL:  "https://media.giphy.com/media/ChFAAc0MssZtC/giphy.gif",
	},
	Motivation{
		Text: "Есть к чему стремиться, за работу",
		URL:  "https://media.giphy.com/media/DXWTC06NVT3QA/giphy.gif",
	},
	Motivation{
		Text: "Сделай подход, может ты не один такой",
		URL:  "https://media.giphy.com/media/T1JzQs6Z7QErC/giphy.gif",
	},
	Motivation{
		Text: "Пора бы подтянуться",
	},
	Motivation{
		Text: "Пора бы подтянуться",
	},
	Motivation{
		Text: "Пора бы подтянуться",
	},
	Motivation{
		Text: "Пора бы подтянуться",
	},
	Motivation{
		Text: "Пора бы подтянуться",
	},
}

func (m *Motivator) SendMotivations() {
	schedules := []Schedule{}
	now := time.Now()
	m.db.Where("at BETWEEN ? AND ?", FromBod(now)-time.Minute*2000, FromBod(now)).Find(&schedules)
	for _, schedule := range schedules {
		var user User
		m.db.Model(&schedule).Related(&user)

		if user.LastWorkout != nil && now.Sub(*user.LastWorkout) < time.Hour*1 {
			continue
		}

		if user.LastWorkout != nil && user.LastScheduleID == &schedule.ID {
			continue
		}
		log.Println(user.Recipient())

		rand.Seed(time.Now().Unix())
		n := rand.Int() % len(motivations)
		m.SendMotivation(&user, motivations[n])

		user.LastWorkout = &now
		user.LastScheduleID = &schedule.ID
		m.db.Save(user)
	}
}

func (m *Motivator) SendMotivation(r tb.Recipient, motivation Motivation) {
	if motivation.URL == "" {
		m.bot.Send(r, motivation.Text)
	}

	m.SendAnimation(r, motivation.URL, motivation.Text)
}

func (m *Motivator) SendAnimation(r tb.Recipient, url string, caption string) {
	motivation := &tb.Animation{
		File:    tb.FromURL(url),
		Caption: caption,
	}
	m.bot.Send(r, motivation)
}
