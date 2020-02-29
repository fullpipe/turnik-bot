package main

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Motivator struct {
	bot *tb.Bot
}

func (m *Motivator) SendAnimation(r tb.Recipient, path string, caption string) {
	motivation := &tb.Animation{
		File:    tb.FromDisk(path),
		Caption: caption,
	}
	m.bot.Send(r, motivation)
}
