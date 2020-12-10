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
		Text: "есть к чему стремиться, за работу",
		URL:  "https://media.giphy.com/media/DXWTC06NVT3QA/giphy.gif",
	},
	Motivation{
		Text: "cделай подход, может ты не один такой",
		URL:  "https://media.giphy.com/media/T1JzQs6Z7QErC/giphy.gif",
	},
	Motivation{Text: "пора бы подтянуться"},
	Motivation{Text: "время повисеть!"},
	Motivation{Text: "еще подходик к перекладине?"},
	Motivation{Text: "ну что, пошли подтянемся!"},
	Motivation{Text: "не забываем подтягиваться!)"},
	Motivation{Text: "турничек тебя совсем заждался!"},
	Motivation{Text: "перекладина соскучилась по тебе!"},
	Motivation{Text: "давай подтянемся!"},
	Motivation{Text: "сам подтянись и позови товарища!"},
	Motivation{Text: "просто сделай это!"},
	Motivation{Text: "выйди во двор, подтянись"},
	Motivation{Text: "ты давно хотел начать, так сделай это сейчас!"},
	Motivation{Text: "погода не оправдание"},
}

var images = []string{
	"https://media.giphy.com/media/l378bt0iIAawyEu0E/giphy.gif",
	"https://media.giphy.com/media/ftegkgLJ5IJzS3b831/giphy.gif",
	"https://media.giphy.com/media/yKE8VoclREyVW/giphy.gif",
	"https://media.giphy.com/media/ftegkgLJ5IJzS3b831/giphy.gif",
	"https://media.giphy.com/media/l378bt0iIAawyEu0E/giphy.gif",
	"https://media.giphy.com/media/MFxTNJG3O16xCKnxQD/giphy.gif",
	"https://media.giphy.com/media/wduVuQiB1QvpC/giphy.gif",
	"https://media.giphy.com/media/LlbcL2bbpRZ7i/giphy.gif",
	"https://media.giphy.com/media/h2ZYTXXR6o7mhMYtQT/giphy.gif",
	"https://media.giphy.com/media/xUPGcD1LxZUkKUMOB2/giphy.gif",
	"https://media.giphy.com/media/IMapVpI6hoS0o/giphy.gif",
	"https://media.giphy.com/media/duxL0JufqG0iQ/giphy.gif",
	"https://media.giphy.com/media/KP19EnxDthURy/giphy.gif",
	"https://media.giphy.com/media/2MKiRiSobn4is/giphy.gif",
	"https://media.giphy.com/media/K4mGIer9AgFNu/giphy.gif",
	"https://media.giphy.com/media/ZrzzYqIgdUaAg/giphy.gif",
	"https://media.giphy.com/media/mTSYOroyZJkNG/giphy.gif",
	"https://media.giphy.com/media/j5VrvBfbV332AOow1B/giphy.gif",
	"https://media.giphy.com/media/zhvsG2MUqW2Ws/giphy.gif",
	"https://media.giphy.com/media/14qd2erccJrq2A/giphy.gif",
	"https://media.giphy.com/media/KZYiIWFAm2X8k/giphy.gif",
	"https://media.giphy.com/media/xT4uQfewy2BixUikIU/giphy.gif",
	"https://media.giphy.com/media/lMwJagHc1tcPe/giphy.gif",
	"https://media.giphy.com/media/CyESeFgx6xgNG/giphy.gif",
	"https://media.giphy.com/media/623ZjXDkiQ9Yk/giphy.gif",
	"https://media.giphy.com/media/gbPGNztPdhS6c/giphy.gif",
	"https://media.giphy.com/media/3o7WTwlO4eajnoKQGk/giphy.gif",
	"https://media.giphy.com/media/iGMpf3IMQxpuw1mHjI/giphy.gif",
	"https://media.giphy.com/media/MuAreXtF0eR9tyKzHR/giphy.gif",
	"https://media.giphy.com/media/Z9tTGDXRjXgI4191cQ/giphy.gif",
	"https://media.giphy.com/media/l1BgRucd74s7erdYs/giphy.gif",
	"https://media.giphy.com/media/HdK972OCf3ahy/giphy.gif",
}

func (m *Motivator) SendMotivations() {
	schedules := []Schedule{}
	now := time.Now()
	m.db.Where("at BETWEEN ? AND ?", FromBod(now)-time.Minute*5, FromBod(now)).Find(&schedules)
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
		go m.SendMotivation(&user, motivations[n])

		user.LastWorkout = &now
		user.LastScheduleID = &schedule.ID
		m.db.Save(user)
	}
}

func (m *Motivator) SendMotivation(r tb.Recipient, motivation Motivation) {
	image := motivation.URL
	if image == "" {
		rand.Seed(time.Now().Unix())
		n := rand.Int() % len(images)
		image = images[n]
	}

	m.SendAnimation(r, image, motivation.Text)
}

func (m *Motivator) SendAnimation(r tb.Recipient, url string, caption string) {
	motivation := &tb.Animation{
		File:    tb.FromURL(url),
		Caption: caption,
	}
	m.bot.Send(r, motivation)
}
