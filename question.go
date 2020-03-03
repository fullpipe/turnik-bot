package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type Question struct {
	Prefix  string
	Text    string
	Bot     *tb.Bot
	answers []Answer
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

	q.answers = []Answer{}
	return q
}

func (q *Question) AddAnswer(id string, text string, callback AnswerCallback) {
	answerButton := tb.InlineButton{
		Unique: q.Prefix + id,
		Text:   text,
	}

	q.Bot.Handle(&answerButton, func(c *tb.Callback) {
		callback(c)
		q.Bot.Respond(c, &tb.CallbackResponse{})
	})

	q.answers = append(q.answers, Answer{
		id:     id,
		text:   text,
		button: &answerButton,
	})
}

func (q *Question) Send(b *tb.Bot, r tb.Recipient, o *tb.SendOptions) (*tb.Message, error) {
	keys := [][]tb.InlineButton{}

	for _, a := range q.answers {
		keys = append(keys, []tb.InlineButton{*a.button})
	}
	return b.Send(r, q.Text, &tb.ReplyMarkup{InlineKeyboard: keys})
}
