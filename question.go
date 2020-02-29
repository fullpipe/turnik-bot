package main

import (
	"log"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Question struct {
	Prefix  string
	Text    string
	Bot     *tb.Bot
	answers map[string]Answer
}

type Answer struct {
	id     string
	text   string
	button *tb.InlineButton
}

type AnswerCallback func(c *tb.Callback)

func NewQuestion(text string, prefix string, bot *tb.Bot) *Question {
	q := &Question{
		Text:   text,
		Prefix: prefix,
		Bot:    bot,
	}

	q.answers = map[string]Answer{}
	return q
}

func (q *Question) AddAnswer(id string, text string, callback AnswerCallback) {
	answerButton := tb.InlineButton{
		Unique: q.Prefix + strconv.Itoa(len(q.answers)),
		Text:   text,
	}
	q.Bot.Handle(&answerButton, func(c *tb.Callback) {
		callback(c)
		q.Bot.Respond(c, &tb.CallbackResponse{})
	})
	q.answers[id] = Answer{
		text:   text,
		button: &answerButton,
	}
}

func (q *Question) Send(b *tb.Bot, r tb.Recipient, o *tb.SendOptions) (*tb.Message, error) {
	keys := [][]tb.InlineButton{}

	for _, a := range q.answers {
		keys = append(keys, []tb.InlineButton{*a.button})
	}
	log.Println("Send question", q.Text)
	return b.Send(r, q.Text, &tb.ReplyMarkup{InlineKeyboard: keys})
}
