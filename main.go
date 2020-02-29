package main

import (
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

var db *gorm.DB
var scheduler *Scheduler

type Schedule struct {
	gorm.Model
	At     time.Duration
	User   User
	UserID uint
}

func main() {
	var err error
	db, err = gorm.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_URI"))

	scheduler = &Scheduler{db: db}

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Schedule{})

	b, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		URL:    os.Getenv("TELEGRAM_URL"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	workTimeQuestion := NewQuestion("Когда ты приходишь на  работу?", "worktime_", b)
	workTimeQuestion.AddAnswer("8", "к 8", func(c *tb.Callback) {
		updateDayStart(c.Sender.Recipient(), "8h")
	})
	workTimeQuestion.AddAnswer("9", "к 9", func(c *tb.Callback) {
		updateDayStart(c.Sender.Recipient(), "9h")
	})
	workTimeQuestion.AddAnswer("10", "к 10", func(c *tb.Callback) {
		updateDayStart(c.Sender.Recipient(), "10h")
	})
	workTimeQuestion.AddAnswer("11", "к 11", func(c *tb.Callback) {
		updateDayStart(c.Sender.Recipient(), "11h")
	})

	howOften := NewQuestion("Как часто ты хотел бы заниматься?", "how_often_", b)
	howOften.AddAnswer("every_hour", "каждый час", func(c *tb.Callback) {
		updateEveryHours(c.Sender.Recipient(), 1)
	})
	howOften.AddAnswer("every_two_hours", "каждые два", func(c *tb.Callback) {
		updateEveryHours(c.Sender.Recipient(), 2)
	})
	howOften.AddAnswer("every_three_hours", "каждые три", func(c *tb.Callback) {
		updateEveryHours(c.Sender.Recipient(), 3)
	})

	c := cron.New()
	sender := &Sender{bot: b, db: db}
	c.AddFunc("@every 5s", func() {
		log.Println("cron")
		sender.SendMotivations()
	})
	c.Start()

	motivator := &Motivator{bot: b}
	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}
		log.Println(m.Sender.Recipient())
		log.Println(m.Sender.ID)

		user := GetOrInitUserById(m.Sender.Recipient())

		if db.NewRecord(user) {
			b.Notify(m.Sender, tb.Typing)
			time.Sleep(2 * time.Second)
			b.Send(m.Sender, "Мы сделаем из тебя человека.")
			b.Notify(m.Sender, tb.UploadingPhoto)
			time.Sleep(2 * time.Second)
			motivator.SendAnimation(m.Sender, "./assets/p2.gif", "Это твоя цель")
			b.Notify(m.Sender, tb.UploadingPhoto)
			time.Sleep(2 * time.Second)
			motivator.SendAnimation(m.Sender, "./assets/b1.gif", "А это ты")

			b.Notify(m.Sender, tb.Typing)
			time.Sleep(2 * time.Second)

			b.Send(m.Sender, "Ответь на пару вопросов")
			db.Save(user)
			b.Notify(m.Sender, tb.Typing)
			time.Sleep(1 * time.Second)
		}

		b.Send(m.Sender, workTimeQuestion)
		b.Notify(m.Sender, tb.Typing)
		time.Sleep(3 * time.Second)
		b.Send(m.Sender, howOften)
	})
	b.Start()
}

func updateDayStart(id string, startTime string) {
	hours, _ := time.ParseDuration(startTime)
	user := GetOrInitUserById(id)
	user.StartsAt = &hours
	db.Save(&user)
	scheduler.Schedule(user)
}

func updateEveryHours(id string, every int) {
	user := GetOrInitUserById(id)
	user.EveryHours = &every
	db.Save(&user)
	scheduler.Schedule(user)
}

func FromBod(t time.Time) time.Duration {
	year, month, day := t.Date()
	return t.Sub(time.Date(year, month, day, 0, 0, 0, 0, t.Location()))
}
