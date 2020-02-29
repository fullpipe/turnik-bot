module github.com/fullpipe/turnik-bot

go 1.13

require (
	github.com/jinzhu/gorm v1.9.12
	github.com/pkg/errors v0.9.1 // indirect
	github.com/robfig/cron/v3 v3.0.1
	gopkg.in/tucnak/telebot.v2 v2.0.0-20200209123123-209b6f88caa9
)

replace gopkg.in/tucnak/telebot.v2 => github.com/fullpipe/telebot v0.0.0-20200229123433-4e97ab59b1ae
